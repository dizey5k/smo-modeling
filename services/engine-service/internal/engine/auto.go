package engine

import (
	"math"
	"smo-gateway/internal/model"
	"smo-gateway/internal/utils"
)

type AutoSimulationResult struct {
	TargetN     int                `json:"target_n"`
	TotalTime   float64            `json:"total_time"`
	SourceStats []model.SourceStat `json:"source_stats"`
	DeviceStats []model.DeviceStat `json:"device_stats"`
}

func RunAutoSimulation(cfg model.SimulationConfig) AutoSimulationResult {
	const (
		tAlpha = 1.643 // alpha = 0.9
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
			// Если отказов 0 или p крайне мало, используем надежную выборку
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

	return finalEngine.CompileStats(nTarget)
}

func (se *SimulationEngine) CompileStats(nTarget int) AutoSimulationResult {
	sourceStats := make([]model.SourceStat, 0, len(se.Config.Sources))
	for _, src := range se.Config.Sources {
		id := src.ID
		gen := se.SourceGenCount[id]
		drop := se.SourceDropCount[id]
		done := se.SourceDoneCount[id]

		pRef := 0.0
		if gen > 0 {
			pRef = float64(drop) / float64(gen)
		}

		mWait, dWait := utils.MeanAndVariance(se.WaitDurations[id])
		mServ, dServ := utils.MeanAndVariance(se.ServDurations[id])
		mStay, _ := utils.MeanAndVariance(se.StayDurations[id])

		sourceStats = append(sourceStats, model.SourceStat{
			SourceID:    id,
			Generated:   gen,
			Dropped:     drop,
			Completed:   done,
			ProbRefusal: pRef,
			TStay:       mStay,
			TWait:       mWait,
			TService:    mServ,
			DispWait:    dWait,
			DispService: dServ,
		})
	}

	devStats := make([]model.DeviceStat, 0, len(se.Pool.Devices))
	for _, dev := range se.Pool.Devices {
		kIsp := 0.0
		if se.CurrentTime > 0 {
			kIsp = dev.BusyDuration / se.CurrentTime
		}
		devStats = append(devStats, model.DeviceStat{
			DeviceID:   dev.ID,
			BusyTime:   dev.BusyDuration,
			UsageCoeff: kIsp,
		})
	}

	return AutoSimulationResult{
		TargetN:     nTarget,
		TotalTime:   se.CurrentTime,
		SourceStats: sourceStats,
		DeviceStats: devStats,
	}
}
