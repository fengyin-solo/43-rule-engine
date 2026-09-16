package handler

import (
	"net/http"

	"ruleengine/pkg/httpx"
)

func (s *Server) registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/health", s.healthCheck)
	mux.HandleFunc("GET /api/health/dependencies", s.healthDependencies)
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	data := s.svc.HealthCheck()
	httpx.OK(w, data)
}

func (s *Server) healthDependencies(w http.ResponseWriter, r *http.Request) {
	data := s.svc.CheckDependencies()
	httpx.OK(w, data)
}
