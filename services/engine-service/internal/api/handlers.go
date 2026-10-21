package api

import (
	"encoding/json"
	"log"
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

	var cfg model.SimulationConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil || len(cfg.Sources) == 0 {
		cfg = h.defaultCfg
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
	log.Printf("[STEP] T=%.4f | Занято слотов: %d | Действие: %s",
		state.CurrentTime, state.TotalGen-state.TotalDone-state.TotalDrop, state.LastAction)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (h *ServerHandler) HandleAuto(w http.ResponseWriter, r *http.Request) {
	var cfg model.SimulationConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil || len(cfg.Sources) == 0 {
		cfg = h.defaultCfg
	}

	results := engine.RunAutoSimulation(cfg)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
