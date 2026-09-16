package handler

import (
	"net/http"

	"ruleengine/internal/model"
	"ruleengine/pkg/httpx"
)

func (s *Server) registerExecutionLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/execution-logs", s.listExecutionLogs)
	mux.HandleFunc("GET /api/execution-logs/{id}", s.getExecutionLog)
	mux.HandleFunc("DELETE /api/execution-logs/{id}", s.deleteExecutionLog)
}

func (s *Server) listExecutionLogs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ExecutionLogFilter{
		RuleSetID: r.URL.Query().Get("rule_set_id"),
		RuleID:    r.URL.Query().Get("rule_id"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	if v := r.URL.Query().Get("hit"); v != "" {
		b := v == "true"
		filter.Hit = &b
	}
	items, total, err := s.svc.ListExecutionLogs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getExecutionLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.GetExecutionLog(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) deleteExecutionLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteExecutionLog(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
