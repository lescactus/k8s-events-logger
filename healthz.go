package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type HealthzServer struct {
	server http.Server
}

type HealthzResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewHealthzServer(addr string) *HealthzServer {
	return &HealthzServer{
		server: http.Server{
			Addr:              addr,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 3 * time.Second,
			WriteTimeout:      5 * time.Second,
		},
	}
}

func (s *HealthzServer) Run() {
	log.Printf("Starting %s http server\n", AppName)
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatalf("Server startup error: %v", err.Error())
	}
}

func (s *HealthzServer) SetRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ready", s.Healthz)
	mux.HandleFunc("/alive", s.Healthz)

	s.server.Handler = mux
}

func (s *HealthzServer) Healthz(w http.ResponseWriter, r *http.Request) {
	healthz := HealthzResponse{Status: http.StatusOK, Message: "Server ready and alive"}

	resp, err := json.Marshal(healthz)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
