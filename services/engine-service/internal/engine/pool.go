package engine

import (
	"math"
	"smo-gateway/internal/model"
	"smo-gateway/internal/utils"
)

type DevicePool struct {
	Devices   []model.Device
	PointerUP int
	MuRate    float64
}

func NewDevicePool(count int, mu float64) *DevicePool {
	devs := make([]model.Device, count)
	for i := 0; i < count; i++ {
		devs[i] = model.Device{ID: i, IsBusy: false}
	}
	return &DevicePool{Devices: devs, PointerUP: 0, MuRate: mu}
}

// SelectDeviceRoundRobin реализует Д2П2 (выбор свободного прибора по кольцу от УП)
func (dp *DevicePool) SelectDeviceRoundRobin() *model.Device {
	n := len(dp.Devices)
	for i := 0; i < n; i++ {
		target := (dp.PointerUP + i) % n
		if !dp.Devices[target].IsBusy {
			dp.PointerUP = (target + 1) % n
			return &dp.Devices[target]
		}
	}
	return nil
}

func (dp *DevicePool) CalcServiceDuration() float64 {
	return utils.Exponential(dp.MuRate)
}

// ExtractTaskByPriority реализует Д2Б4 (И1 > И2 > И3, при равенстве - LIFO)
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
