package engine

import (
	"testing"

	"smo-gateway/internal/model"
)

// Тест Д1ОЗ1: циклическая запись и сдвиг УБ
func TestRingBuffer_D1OZ1(t *testing.T) {
	buf := NewRingBuffer(4)

	// Записываем 3 задачи
	for i := 1; i <= 3; i++ {
		task := &model.InferenceTask{ID: i, SourceID: 1, Priority: 1}
		dropped, slot := buf.PutTask(task, float64(i))

		if dropped != nil {
			t.Fatalf("Ожидалось отсутствие отказа, но выбита задача #%d", dropped.ID)
		}
		if slot != i-1 {
			t.Errorf("Задача #%d должна быть записана в слот %d, но попала в %d", i, i-1, slot)
		}
		if buf.GetPointer() != i%4 {
			t.Errorf("УБ должен быть равен %d, получено %d", i%4, buf.GetPointer())
		}
	}
}

// Тест Д1ОО1: выбивание заявки строго под указателем УБ при переполнении буфера
func TestRingBuffer_D1OO1_Drop(t *testing.T) {
	cap := 3
	buf := NewRingBuffer(cap)

	// Заполняем буфер до предела (слоты 0, 1, 2). УБ вернется в 0.
	for i := 1; i <= cap; i++ {
		buf.PutTask(&model.InferenceTask{ID: i, SourceID: 1, Priority: 1}, float64(i))
	}

	if !buf.IsFull() {
		t.Fatal("Буфер должен быть полон")
	}
	if buf.GetPointer() != 0 {
		t.Fatalf("УБ должен указывать на слот 0, текущий: %d", buf.GetPointer())
	}

	// Попытка вставки 4-й задачи. Должна быть выбита задача #1 из слота 0 под УБ
	taskNew := &model.InferenceTask{ID: 4, SourceID: 2, Priority: 2}
	dropped, slot := buf.PutTask(taskNew, 10.0)

	if dropped == nil {
		t.Fatal("Ожидался сброс заявки по правилу Д1ОО1, но сброса не произошло")
	}
	if dropped.ID != 1 {
		t.Errorf("Ожидался сброс задачи #1, сброшена #%d", dropped.ID)
	}
	if slot != 0 {
		t.Errorf("Новая задача должна занять слот 0, занят %d", slot)
	}
	if buf.GetPointer() != 1 {
		t.Errorf("После выбивания УБ должен сдвинуться на 1, текущий: %d", buf.GetPointer())
	}
}

// Тест Д2Б4: приоритетный выбор заявок (И1 > И2 > И3) + LIFO для одного источника
func TestRingBuffer_D2B4_Priority(t *testing.T) {
	buf := NewRingBuffer(5)

	// Добавляем задачи в разнобой по приоритетам и времени
	buf.PutTask(&model.InferenceTask{ID: 101, SourceID: 3, Priority: 3, GenerationTime: 1.0}, 1.0)
	buf.PutTask(&model.InferenceTask{ID: 102, SourceID: 2, Priority: 2, GenerationTime: 2.0}, 2.0)
	buf.PutTask(&model.InferenceTask{ID: 103, SourceID: 1, Priority: 1, GenerationTime: 3.0}, 3.0) // Высший
	buf.PutTask(&model.InferenceTask{ID: 104, SourceID: 2, Priority: 2, GenerationTime: 4.0}, 4.0) // LIFO для И2

	// 1. Первая должна извлечься задача И1 (#103)
	t1 := buf.ExtractTaskByPriority()
	if t1 == nil || t1.ID != 103 {
		t.Fatalf("Ожидалась задача #103 (И1, высший приоритет), получена: %+v", t1)
	}

	// 2. Вторая — задача И2 с более свежим временем (#104, LIFO)
	t2 := buf.ExtractTaskByPriority()
	if t2 == nil || t2.ID != 104 {
		t.Fatalf("Ожидалась задача #104 (И2, правило LIFO), получена: %+v", t2)
	}

	// 3. Третья — оставшаяся задача И2 (#102)
	t3 := buf.ExtractTaskByPriority()
	if t3 == nil || t3.ID != 102 {
		t.Fatalf("Ожидалась задача #102, получена: %+v", t3)
	}

	// 4. Последняя — низкоприоритетная И3 (#101)
	t4 := buf.ExtractTaskByPriority()
	if t4 == nil || t4.ID != 101 {
		t.Fatalf("Ожидалась задача #101, получена: %+v", t4)
	}
}

// Тест Д2П2: циклический поиск свободного прибора с указателем УП
func TestDevicePool_D2P2_RoundRobin(t *testing.T) {
	pool := NewDevicePool(3, 1.0)

	// На старте УП = 0 -> должен выбрать прибор 0, УП станет 1
	d1 := pool.SelectDeviceRoundRobin()
	if d1.ID != 0 || pool.PointerUP != 1 {
		t.Errorf("Ожидался прибор 0 и УП=1, получено: ID=%d, УП=%d", d1.ID, pool.PointerUP)
	}
	d1.IsBusy = true

	// Следующий выбор: свободен прибор 1, УП станет 2
	d2 := pool.SelectDeviceRoundRobin()
	if d2.ID != 1 || pool.PointerUP != 2 {
		t.Errorf("Ожидался прибор 1 и УП=2, получено: ID=%d, УП=%d", d2.ID, pool.PointerUP)
	}
	d2.IsBusy = true

	// Следующий выбор: свободен прибор 2, УП станет (2+1)%3 = 0
	d3 := pool.SelectDeviceRoundRobin()
	if d3.ID != 2 || pool.PointerUP != 0 {
		t.Errorf("Ожидался прибор 2 и УП=0, получено: ID=%d, УП=%d", d3.ID, pool.PointerUP)
	}
	d3.IsBusy = true

	// Все приборы заняты -> nil
	if pool.SelectDeviceRoundRobin() != nil {
		t.Error("Ожидался nil, так как все приборы заняты")
	}

	// Освобождаем прибор 1 (при УП=0). Поиск от УП=0 найдет свободный прибор 1 и сдвинет УП на 2
	pool.Devices[1].IsBusy = false
	dNext := pool.SelectDeviceRoundRobin()
	if dNext.ID != 1 || pool.PointerUP != 2 {
		t.Errorf("Ожидался прибор 1 и УП=2, получено: ID=%d, УП=%d", dNext.ID, pool.PointerUP)
	}
}

// Тест изоляции Зонного буфера: переполнение одного источника не выбивает заявки других источников
func TestZonedBuffer_Isolation(t *testing.T) {
	// Буфер на 6 ячеек для 3 источников (по 2 ячейки на зону)
	zb := NewZonedBuffer(6, 3)

	// Источник 1 занимает свои 2 ячейки
	zb.PutTask(&model.InferenceTask{ID: 1, SourceID: 1, Priority: 1}, 1.0)
	zb.PutTask(&model.InferenceTask{ID: 2, SourceID: 1, Priority: 1}, 2.0)

	// Источник 2 занимает 1 ячейку
	zb.PutTask(&model.InferenceTask{ID: 10, SourceID: 2, Priority: 2}, 3.0)

	// Третья задача от Источника 1 переполняет зону 1.
	// Должна выбиться задача #1 из зоны 1, а задача #10 от Источника 2 обязана остаться нетронутой!
	dropped, _ := zb.PutTask(&model.InferenceTask{ID: 3, SourceID: 1, Priority: 1}, 4.0)

	if dropped == nil || dropped.ID != 1 {
		t.Fatalf("Ожидался сброс задачи #1 из зоны источника 1, получено: %+v", dropped)
	}

	slots := zb.GetSlots()
	hasTask10 := false
	for _, slot := range slots {
		if slot.IsOccupied && slot.Task != nil && slot.Task.ID == 10 {
			hasTask10 = true
			break
		}
	}

	if !hasTask10 {
		t.Fatal("Ошибка изоляции: задача источника 2 была повреждена или вытеснена чужим переполнением!")
	}
}
