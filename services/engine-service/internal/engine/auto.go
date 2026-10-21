package engine

import (
	"math"
	"smo-gateway/internal/model"
)

type AutoSimulationResult struct {
	TargetN     int                `json:"target_n"`
	TotalTime   float64            `json:"total_time"`
	BufferMode  string             `json:"buffer_mode"`
	SourceStats []model.SourceStat `json:"source_stats"`
	DeviceStats []model.DeviceStat `json:"device_stats"`
}

// CompileStats собирает финальную статистику СМО по источникам и приборам (ОР1)
func (se *SimulationEngine) CompileStats(targetN int) AutoSimulationResult {
	sourceStats := make([]model.SourceStat, 0, len(se.Config.Sources))

	for _, src := range se.Config.Sources {
		gen := se.SourceGenCount[src.ID]
		drop := se.SourceDropCount[src.ID]
		done := se.SourceDoneCount[src.ID]

		potk := 0.0
		if gen > 0 {
			potk = float64(drop) / float64(gen)
		}

		tWait, dWait := calcMeanAndVariance(se.WaitDurations[src.ID])
		tServ, dServ := calcMeanAndVariance(se.ServDurations[src.ID])
		tStay, _ := calcMeanAndVariance(se.StayDurations[src.ID])

		sourceStats = append(sourceStats, model.SourceStat{
			SourceID:  src.ID,
			Generated: gen,
			Dropped:   drop,
			Completed: done,
			POtk:      potk,
			TStay:     tStay,
			TWait:     tWait,
			TServ:     tServ,
			DWait:     dWait,
			DServ:     dServ,
		})
	}

	deviceStats := make([]model.DeviceStat, 0, len(se.Pool.Devices))
	for _, dev := range se.Pool.Devices {
		kIsp := 0.0
		if se.CurrentTime > 0 {
			kIsp = dev.BusyDuration / se.CurrentTime
			if kIsp > 1.0 {
				kIsp = 1.0
			}
		}

		deviceStats = append(deviceStats, model.DeviceStat{
			DeviceID: dev.ID,
			BusyTime: dev.BusyDuration,
			KIsp:     kIsp,
		})
	}

	return AutoSimulationResult{
		TargetN:     targetN,
		TotalTime:   se.CurrentTime,
		BufferMode:  se.Config.BufferMode,
		SourceStats: sourceStats,
		DeviceStats: deviceStats,
	}
}

func calcMeanAndVariance(vals []float64) (mean float64, variance float64) {
	n := len(vals)
	if n == 0 {
		return 0, 0
	}

	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	mean = sum / float64(n)

	sumSq := 0.0
	for _, v := range vals {
		diff := v - mean
		sumSq += diff * diff
	}
	variance = sumSq / float64(n)
	return mean, variance
}

// RunAutoSimulation запускает автоматический прогон с калибровкой выборки по правилу 10%
func RunAutoSimulation(cfg model.SimulationConfig) AutoSimulationResult {
	const (
		tAlpha = 1.643 // alpha = 0.9 для доверительного интервала
		delta  = 0.10  // относительная погрешность 10%
	)

	nTarget := 100
	var finalEngine *SimulationEngine

	for {
		engine := NewSimulationEngine(cfg)
		for engine.TotalGen < nTarget && len(engine.Calendar) > 0 {
			engine.Step()
		}

		pGlobal := 0.0
		if engine.TotalGen > 0 {
			pGlobal = float64(engine.TotalDrop) / float64(engine.TotalGen)
		}

		var nextN int
		if pGlobal > 0.001 && pGlobal < 0.999 {
			nextN = int(math.Ceil((math.Pow(tAlpha, 2) * (1.0 - pGlobal)) / (pGlobal * math.Pow(delta, 2))))
			if nextN < 100 {
				nextN = 100
			}
			if nextN > 50000 {
				nextN = 50000
			}
		} else {
			nextN = 3000
		}

		if nTarget >= nextN || math.Abs(float64(nextN-nTarget))/float64(nTarget) < 0.10 {
			finalEngine = engine
			break
		}
		nTarget = nextN
	}

	for !finalEngine.Buffer.IsEmpty() && len(finalEngine.Calendar) > 0 {
		finalEngine.Step()
	}

	res := finalEngine.CompileStats(nTarget)
	res.BufferMode = cfg.BufferMode
	return res
}
