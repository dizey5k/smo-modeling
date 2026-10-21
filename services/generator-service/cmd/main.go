package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"smo-generator/internal/distr"
)

type BatchRequest struct {
	Count  int     `json:"count"`
	Type   string  `json:"type"` // "UNIFORM" | "EXPONENTIAL"
	RangeA float64 `json:"range_a,omitempty"`
	RangeB float64 `json:"range_b,omitempty"`
	Lambda float64 `json:"lambda,omitempty"`
}

type BatchResponse struct {
	Values []float64 `json:"values"`
}

func handleUniform(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	a, _ := strconv.ParseFloat(q.Get("a"), 64)
	b, _ := strconv.ParseFloat(q.Get("b"), 64)
	if b <= a {
		b = a + 1.0
	}

	val := distr.Uniform(a, b)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"value": val})
}

func handleExponential(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	lambda, _ := strconv.ParseFloat(q.Get("lambda"), 64)
	if lambda <= 0 {
		lambda = 1.0
	}

	val := distr.Exponential(lambda)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"value": val})
}

// Пакетная генерация для авто-режима (чтобы не спамить тысячами HTTP-запросов)
func handleBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Count <= 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	values := make([]float64, req.Count)
	if req.Type == "UNIFORM" {
		for i := 0; i < req.Count; i++ {
			values[i] = distr.Uniform(req.RangeA, req.RangeB)
		}
	} else {
		for i := 0; i < req.Count; i++ {
			values[i] = distr.Exponential(req.Lambda)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(BatchResponse{Values: values})
}

func main() {
	http.HandleFunc("/api/distr/uniform", handleUniform)
	http.HandleFunc("/api/distr/exponential", handleExponential)
	http.HandleFunc("/api/distr/batch", handleBatch)

	log.Println("[GENERATOR-SERVICE] Слушает на :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
