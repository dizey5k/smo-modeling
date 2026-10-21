package engine

import (
	"math"
	"smo-gateway/internal/model"
)

type ZonedBuffer struct {
	capacity    int
	sourceCount int
	zones       map[int][]int // sourceID -> индексы слотов
	pointers    map[int]int   // sourceID -> локальный указатель зоны
	slots       []model.BufferSlot
}

func NewZonedBuffer(totalCap, sourcesCount int) *ZonedBuffer {
	slots := make([]model.BufferSlot, totalCap)
	for i := 0; i < totalCap; i++ {
		slots[i] = model.BufferSlot{Index: i, IsOccupied: false, Task: nil}
	}

	zones := make(map[int][]int)
	pointers := make(map[int]int)

	slotsPerZone := totalCap / sourcesCount
	extra := totalCap % sourcesCount

	currentIdx := 0
	for s := 1; s <= sourcesCount; s++ {
		count := slotsPerZone
		if s <= extra {
			count++
		}
		zoneSlots := make([]int, 0, count)
		for j := 0; j < count; j++ {
			zoneSlots = append(zoneSlots, currentIdx)
			currentIdx++
		}
		zones[s] = zoneSlots
		pointers[s] = 0
	}

	return &ZonedBuffer{
		capacity:    totalCap,
		sourceCount: sourcesCount,
		zones:       zones,
		pointers:    pointers,
		slots:       slots,
	}
}

func (zb *ZonedBuffer) GetSlots() []model.BufferSlot {
	res := make([]model.BufferSlot, len(zb.slots))
	copy(res, zb.slots)
	return res
}

func (zb *ZonedBuffer) GetPointer() int {
	return 0
}

func (zb *ZonedBuffer) IsFull() bool {
	for _, slot := range zb.slots {
		if !slot.IsOccupied {
			return false
		}
	}
	return true
}

func (zb *ZonedBuffer) IsEmpty() bool {
	for _, slot := range zb.slots {
		if slot.IsOccupied {
			return false
		}
	}
	return true
}

func (zb *ZonedBuffer) PutTask(task *model.InferenceTask, currentTime float64) (*model.InferenceTask, int) {
	task.BufferEntryTime = currentTime
	task.Status = model.StatusInBuffer

	zoneSlots, exists := zb.zones[task.SourceID]
	if !exists || len(zoneSlots) == 0 {
		task.Status = model.StatusDropped
		return task, -1
	}

	zoneLen := len(zoneSlots)
	localPtr := zb.pointers[task.SourceID]

	// 1. Поиск свободного слота внутри зоны источника по кольцу
	for i := 0; i < zoneLen; i++ {
		slotIdx := zoneSlots[(localPtr+i)%zoneLen]
		if !zb.slots[slotIdx].IsOccupied {
			zb.slots[slotIdx].IsOccupied = true
			zb.slots[slotIdx].Task = task
			zb.pointers[task.SourceID] = (localPtr + i + 1) % zoneLen
			return nil, slotIdx
		}
	}

	// 2. Зона переполнена -> выбиваем заявку под указателем зоны (Д1ОО1 в рамках зоны)
	droppedSlotIdx := zoneSlots[localPtr]
	droppedTask := zb.slots[droppedSlotIdx].Task
	droppedTask.Status = model.StatusDropped

	zb.slots[droppedSlotIdx].Task = task
	zb.pointers[task.SourceID] = (localPtr + 1) % zoneLen

	return droppedTask, droppedSlotIdx
}

func (zb *ZonedBuffer) ExtractTaskByPriority() *model.InferenceTask {
	bestSlot := -1
	bestPriority := math.MaxInt32
	var latestGenTime float64 = -1.0

	for _, slot := range zb.slots {
		if !slot.IsOccupied || slot.Task == nil {
			continue
		}
		t := slot.Task
		if t.Priority < bestPriority {
			bestPriority = t.Priority
			bestSlot = slot.Index
			latestGenTime = t.GenerationTime
		} else if t.Priority == bestPriority && t.GenerationTime > latestGenTime {
			bestSlot = slot.Index
			latestGenTime = t.GenerationTime
		}
	}

	if bestSlot == -1 {
		return nil
	}

	task := zb.slots[bestSlot].Task
	zb.slots[bestSlot].IsOccupied = false
	zb.slots[bestSlot].Task = nil
	return task
}
