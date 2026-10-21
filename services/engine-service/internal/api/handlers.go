package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"smo-gateway/internal/engine"
	"smo-gateway/internal/model"
)

type ServerHandler struct {
	mu         sync.Mutex
	engine     *engine.SimulationEngine
	defaultCfg model.SimulationConfig
}

func NewServerHandler(cfg model.SimulationConfig) *ServerHandler {
	return &ServerHandler{defaultCfg: cfg}
}

func (h *ServerHandler) HandleInit(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	cfg := h.defaultCfg
	if r.Body != nil {
		var userCfg model.SimulationConfig
		if err := json.NewDecoder(r.Body).Decode(&userCfg); err == nil {
			if len(userCfg.Sources) > 0 {
				cfg = userCfg
			} else if userCfg.BufferMode != "" {
				cfg.BufferMode = userCfg.BufferMode
			}
		}
	}

	h.engine = engine.NewSimulationEngine(cfg)
	state := h.engine.GetState("Система инициализирована", false)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (h *ServerHandler) HandleStep(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.engine == nil {
		h.engine = engine.NewSimulationEngine(h.defaultCfg)
	}

	state := h.engine.Step()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (h *ServerHandler) HandleAuto(w http.ResponseWriter, r *http.Request) {
	cfg := h.defaultCfg
	if r.Body != nil {
		var userCfg model.SimulationConfig
		if err := json.NewDecoder(r.Body).Decode(&userCfg); err == nil {
			if len(userCfg.Sources) > 0 {
				cfg = userCfg
			} else if userCfg.BufferMode != "" {
				cfg.BufferMode = userCfg.BufferMode
			}
		}
	}

	results := engine.RunAutoSimulation(cfg)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *ServerHandler) HandleExperiment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		model.SimulationConfig
		Experiment model.ExperimentRequest `json:"experiment"`
	}

	cfg := h.defaultCfg
	expReq := model.ExperimentRequest{}

	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			if len(req.Sources) > 0 {
				cfg = req.SimulationConfig
			}
			expReq = req.Experiment
		}
	}

	results := engine.RunGridSearch(cfg, expReq)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
