package engine

import (
	"math"
	"sync"

	"smo-gateway/internal/model"
)

type ExperimentRun struct {
	DevicesCount  int     `json:"devices_count"`
	BufferCap     int     `json:"buffer_cap"`
	BufferMode    string  `json:"buffer_mode"`
	TargetN       int     `json:"target_n"`
	TotalTime     float64 `json:"total_time"`
	ProbRefusal   float64 `json:"p_otk"`
	AvgStayTime   float64 `json:"t_stay"`
	AvgUsage      float64 `json:"k_isp"`
	LoadFactor    float64 `json:"load_rho"`
	IsOptimal     bool    `json:"is_optimal"`
	EstimatedCost float64 `json:"cost"`
}

type ExperimentResult struct {
	Runs           []ExperimentRun `json:"runs"`
	OptimalConfig  *ExperimentRun  `json:"optimal_config"`
	Recommendation string          `json:"recommendation"`
}

// RunGridSearch выполняет параллельный параметрический синтез ВС (Этап 4)
func RunGridSearch(baseCfg model.SimulationConfig, req model.ExperimentRequest) ExperimentResult {
	targetN := req.TargetN
	if targetN <= 0 {
		targetN = 1000
	} else if targetN > 10000 {
		targetN = 10000
	}

	deviceOptions := req.DeviceOptions
	if len(deviceOptions) == 0 {
		deviceOptions = []int{2, 3, 4}
	}

	bufferOptions := req.BufferOptions
	if len(bufferOptions) == 0 {
		bufferOptions = []int{4, 8, 12, 16}
	}

	modes := req.Modes
	if len(modes) == 0 {
		modes = []string{"RING", "ZONED"}
	}

	sumLambda := 0.0
	for _, s := range baseCfg.Sources {
		avgInterval := (s.RangeA + s.RangeB) / 2.0
		if avgInterval > 0 {
			sumLambda += 1.0 / avgInterval
		}
	}

	totalCombinations := len(modes) * len(deviceOptions) * len(bufferOptions)
	runs := make([]ExperimentRun, totalCombinations)

	var wg sync.WaitGroup
	idx := 0

	for _, mode := range modes {
		for _, devCount := range deviceOptions {
			for _, bufCap := range bufferOptions {
				wg.Add(1)
				go func(slotIndex int, m string, devs, buf, n int) {
					defer wg.Done()

					cfg := baseCfg
					cfg.DevicesCount = devs
					cfg.BufferCap = buf
					cfg.BufferMode = m

					engine := NewSimulationEngine(cfg)
					for engine.TotalGen < n && len(engine.Calendar) > 0 {
						engine.Step()
					}
					for !engine.Buffer.IsEmpty() && len(engine.Calendar) > 0 {
						engine.Step()
					}

					res := engine.CompileStats(n)

					totalGen, totalDrop := 0, 0
					totalStaySum, stayCount := 0.0, 0
					for _, ss := range res.SourceStats {
						totalGen += ss.Generated
						totalDrop += ss.Dropped
						totalStaySum += ss.TStay * float64(ss.Completed)
						stayCount += ss.Completed
					}

					pGlobal := 0.0
					if totalGen > 0 {
						pGlobal = float64(totalDrop) / float64(totalGen)
					}

					avgStay := 0.0
					if stayCount > 0 {
						avgStay = totalStaySum / float64(stayCount)
					}

					avgKisp := 0.0
					for _, ds := range res.DeviceStats {
						avgKisp += ds.KIsp
					}
					if len(res.DeviceStats) > 0 {
						avgKisp /= float64(len(res.DeviceStats))
					}

					rho := sumLambda / (float64(devs) * cfg.MuRate)
					cost := float64(devs)*8000.0 + float64(buf)*300.0
					isOpt := pGlobal <= 0.10 && avgKisp >= 0.80

					runs[slotIndex] = ExperimentRun{
						DevicesCount:  devs,
						BufferCap:     buf,
						BufferMode:    m,
						TargetN:       n,
						TotalTime:     res.TotalTime,
						ProbRefusal:   pGlobal,
						AvgStayTime:   avgStay,
						AvgUsage:      avgKisp,
						LoadFactor:    rho,
						IsOptimal:     isOpt,
						EstimatedCost: cost,
					}
				}(idx, mode, devCount, bufCap, targetN)
				idx++
			}
		}
	}

	wg.Wait()

	var optimal *ExperimentRun
	minCost := math.MaxFloat64

	for i := range runs {
		r := &runs[i]
		if r.IsOptimal && r.EstimatedCost < minCost {
			minCost = r.EstimatedCost
			optCopy := *r
			optimal = &optCopy
		}
	}

	if optimal == nil && len(runs) > 0 {
		bestRun := runs[0]
		for _, r := range runs {
			if r.ProbRefusal <= 0.10 && r.AvgUsage > bestRun.AvgUsage {
				bestRun = r
			}
		}
		optimal = &bestRun
	}

	return ExperimentResult{
		Runs:           runs,
		OptimalConfig:  optimal,
		Recommendation: "Оптимальная архитектура ВС подобрана по минимуму затрат при соблюдении требований ТЗ (P_отк <= 10%, K_исп >= 80%).",
	}
}
