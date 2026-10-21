package engine

import (
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
