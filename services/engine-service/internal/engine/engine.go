package engine

import (
	"sort"

	"smo-gateway/internal/client"
	"smo-gateway/internal/model"
)

type SimulationEngine struct {
	Config      model.SimulationConfig
	CurrentTime float64
	Calendar    []model.Event
	Buffer      *RingBuffer
	Pool        *DevicePool
	GenClient   *client.GeneratorClient
	Tasks       map[int]*model.InferenceTask
	TaskSeq     int

	TotalGen  int
	TotalDrop int
	TotalDone int

	SourceGenCount  map[int]int
	SourceDropCount map[int]int
	SourceDoneCount map[int]int
	WaitDurations   map[int][]float64
	ServDurations   map[int][]float64
	StayDurations   map[int][]float64
}

func NewSimulationEngine(cfg model.SimulationConfig) *SimulationEngine {
	genClient := client.NewGeneratorClient()

	se := &SimulationEngine{
		Config:          cfg,
		Buffer:          NewRingBuffer(cfg.BufferCap),
		Pool:            NewDevicePool(cfg.DevicesCount, cfg.MuRate),
		GenClient:       genClient,
		Tasks:           make(map[int]*model.InferenceTask),
		SourceGenCount:  make(map[int]int),
		SourceDropCount: make(map[int]int),
		SourceDoneCount: make(map[int]int),
		WaitDurations:   make(map[int][]float64),
		ServDurations:   make(map[int][]float64),
		StayDurations:   make(map[int][]float64),
	}

	// ИЗ2: равномерный закон
	for _, src := range cfg.Sources {
		interval := se.GenClient.GetUniform(src.RangeA, src.RangeB)
		se.ScheduleEvent(model.Event{
			Time:     interval,
			Type:     model.EventTaskGenerated,
			SourceID: src.ID,
		})
	}
	return se
}

func (se *SimulationEngine) ScheduleEvent(e model.Event) {
	se.Calendar = append(se.Calendar, e)
	sort.Slice(se.Calendar, func(i, j int) bool {
		return se.Calendar[i].Time < se.Calendar[j].Time
	})
}

func (se *SimulationEngine) GetState(action string, finished bool) model.SimulationState {
	devsCopy := make([]model.Device, len(se.Pool.Devices))
	copy(devsCopy, se.Pool.Devices)

	eventsCopy := make([]model.Event, len(se.Calendar))
	copy(eventsCopy, se.Calendar)

	return model.SimulationState{
		CurrentTime: se.CurrentTime,
		PointerUB:   se.Buffer.GetPointer(),
		PointerUP:   se.Pool.PointerUP,
		Buffer:      se.Buffer.GetSlots(),
		Devices:     devsCopy,
		Events:      eventsCopy,
		LastAction:  action,
		TotalGen:    se.TotalGen,
		TotalDrop:   se.TotalDrop,
		TotalDone:   se.TotalDone,
		IsFinished:  finished,
	}
}
