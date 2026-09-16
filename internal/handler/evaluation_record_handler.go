package handler

import (
	"net/http"

	"ruleengine/internal/model"
	"ruleengine/pkg/httpx"
)

func (s *Server) registerEvaluationRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/evaluation-records", s.listEvaluationRecords)
	mux.HandleFunc("GET /api/evaluation-records/{id}", s.getEvaluationRecord)
	mux.HandleFunc("DELETE /api/evaluation-records/{id}", s.deleteEvaluationRecord)
}

func (s *Server) listEvaluationRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.EvaluationRecordFilter{
		RuleSetID: r.URL.Query().Get("rule_set_id"),
		Result:    r.URL.Query().Get("result"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListEvaluationRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getEvaluationRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.GetEvaluationRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) deleteEvaluationRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteEvaluationRecord(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
