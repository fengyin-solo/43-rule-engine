package handler

import (
	"net/http"
	"strconv"
	"time"

	"ruleengine/pkg/httpx"
)

func (s *Server) registerCleanupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/cleanup/evaluation-records", s.cleanupEvaluations)
	mux.HandleFunc("POST /api/cleanup/execution-logs", s.cleanupLogs)
	mux.HandleFunc("POST /api/cleanup/audit-records", s.cleanupAudits)
	mux.HandleFunc("POST /api/cleanup/rule-versions", s.cleanupRuleVersions)
}

func (s *Server) cleanupEvaluations(w http.ResponseWriter, r *http.Request) {
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days < 1 {
		days = 30
	}
	before := time.Now().AddDate(0, 0, -days)
	count, err := s.svc.CleanupOldEvaluationRecords(before)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"deleted": count})
}

func (s *Server) cleanupLogs(w http.ResponseWriter, r *http.Request) {
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days < 1 {
		days = 30
	}
	before := time.Now().AddDate(0, 0, -days)
	count, err := s.svc.CleanupOldExecutionLogs(before)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"deleted": count})
}

func (s *Server) cleanupAudits(w http.ResponseWriter, r *http.Request) {
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days < 1 {
		days = 90
	}
	before := time.Now().AddDate(0, 0, -days)
	count, err := s.svc.CleanupOldAuditRecords(before)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"deleted": count})
}

func (s *Server) cleanupRuleVersions(w http.ResponseWriter, r *http.Request) {
	keep, _ := strconv.Atoi(r.URL.Query().Get("keep"))
	count, err := s.svc.CleanupOldRuleVersions(keep)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"deleted": count})
}
