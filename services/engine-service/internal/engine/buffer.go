package engine

import "smo-gateway/internal/model"

type RingBuffer struct {
	capacity  int
	pointerUB int
	slots     []model.BufferSlot
}

func NewRingBuffer(cap int) *RingBuffer {
	slots := make([]model.BufferSlot, cap)
	for i := 0; i < cap; i++ {
		slots[i] = model.BufferSlot{Index: i, IsOccupied: false, Task: nil}
	}
	return &RingBuffer{
		capacity:  cap,
		pointerUB: 0,
		slots:     slots,
	}
}

func (rb *RingBuffer) GetPointer() int {
	return rb.pointerUB
}

func (rb *RingBuffer) GetSlots() []model.BufferSlot {
	res := make([]model.BufferSlot, len(rb.slots))
	copy(res, rb.slots)
	return res
}

func (rb *RingBuffer) IsFull() bool {
	for _, slot := range rb.slots {
		if !slot.IsOccupied {
			return false
		}
	}
	return true
}

func (rb *RingBuffer) IsEmpty() bool {
	for _, slot := range rb.slots {
		if slot.IsOccupied {
			return false
		}
	}
	return true
}

// PutTask реализует Д1ОЗ1 (постановка по кольцу) и Д1ОО1 (отказ под указателем)
func (rb *RingBuffer) PutTask(task *model.InferenceTask, currentTime float64) (droppedTask *model.InferenceTask, slotIdx int) {
	task.BufferEntryTime = currentTime
	task.Status = model.StatusInBuffer

	// 1. Поиск свободного слота по кольцу начиная с указателя УБ (Д1ОЗ1)
	for i := 0; i < rb.capacity; i++ {
		target := (rb.pointerUB + i) % rb.capacity
		if !rb.slots[target].IsOccupied {
			rb.slots[target].IsOccupied = true
			rb.slots[target].Task = task
			rb.pointerUB = (target + 1) % rb.capacity
			return nil, target
		}
	}

	// 2. Буфер полон -> Срабатывает Д1ОО1 (выбивание под указателем)
	droppedSlot := rb.pointerUB
	droppedTask = rb.slots[droppedSlot].Task
	droppedTask.Status = model.StatusDropped

	// Перезаписываем слот новой заявкой
	rb.slots[droppedSlot].Task = task
	rb.pointerUB = (droppedSlot + 1) % rb.capacity

	return droppedTask, droppedSlot
}
