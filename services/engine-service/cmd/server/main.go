package main

import (
	"log"
	"net/http"

	"smo-gateway/internal/api"
	"smo-gateway/internal/model"
)

func main() {
	defaultConfig := model.SimulationConfig{
		SourcesCount: 3,
		BufferCap:    8,
		DevicesCount: 3,
		MuRate:       1.2,
		BufferMode:   "RING",
		Sources: []model.SourceConfig{
			{ID: 1, Priority: 1, RangeA: 0.5, RangeB: 1.5},
			{ID: 2, Priority: 2, RangeA: 1.0, RangeB: 2.0},
			{ID: 3, Priority: 3, RangeA: 1.5, RangeB: 3.0},
		},
	}

	handler := api.NewServerHandler(defaultConfig)

	http.HandleFunc("/api/init", api.WithLogger(handler.HandleInit))
	http.HandleFunc("/api/step", api.WithLogger(handler.HandleStep))
	http.HandleFunc("/api/auto", api.WithLogger(handler.HandleAuto))

	log.Println("SMO Engine Service started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
