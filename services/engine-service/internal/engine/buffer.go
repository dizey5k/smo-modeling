package engine

import (
	"math"
	"smo-gateway/internal/model"
)

// BufferInterface задаёт контракт для кольцевого и зонного буфера
type BufferInterface interface {
	PutTask(task *model.InferenceTask, currentTime float64) (*model.InferenceTask, int)
	ExtractTaskByPriority() *model.InferenceTask
	GetSlots() []model.BufferSlot
	GetPointer() int
	IsFull() bool
	IsEmpty() bool
}

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
func (rb *RingBuffer) PutTask(task *model.InferenceTask, currentTime float64) (*model.InferenceTask, int) {
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
	droppedTask := rb.slots[droppedSlot].Task
	droppedTask.Status = model.StatusDropped

	rb.slots[droppedSlot].Task = task
	rb.pointerUB = (droppedSlot + 1) % rb.capacity

	return droppedTask, droppedSlot
}

// ExtractTaskByPriority реализует Д2Б4 (приоритет И1 > И2 > И3, при равенстве — LIFO)
func (rb *RingBuffer) ExtractTaskByPriority() *model.InferenceTask {
	bestSlot := -1
	bestPriority := math.MaxInt32
	var latestGenTime float64 = -1.0

	for i := 0; i < rb.capacity; i++ {
		if !rb.slots[i].IsOccupied || rb.slots[i].Task == nil {
			continue
		}
		t := rb.slots[i].Task
		if t.Priority < bestPriority {
			bestPriority = t.Priority
			bestSlot = i
			latestGenTime = t.GenerationTime
		} else if t.Priority == bestPriority && t.GenerationTime > latestGenTime {
			bestSlot = i
			latestGenTime = t.GenerationTime
		}
	}

	if bestSlot == -1 {
		return nil
	}

	task := rb.slots[bestSlot].Task
	rb.slots[bestSlot].IsOccupied = false
	rb.slots[bestSlot].Task = nil
	return task
}
