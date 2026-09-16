package handler

import (
	"net/http"

	"ruleengine/internal/model"
	"ruleengine/pkg/httpx"
)

func (s *Server) registerAuditRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/audit-records", s.createAuditRecord)
	mux.HandleFunc("GET /api/audit-records", s.listAuditRecords)
	mux.HandleFunc("GET /api/audit-records/{id}", s.getAuditRecord)
	mux.HandleFunc("DELETE /api/audit-records/{id}", s.deleteAuditRecord)
}

type createAuditRecordRequest struct {
	Operator       string `json:"operator"`
	Operation      string `json:"operation"`
	TargetType     string `json:"target_type"`
	TargetID       string `json:"target_id"`
	BeforeSnapshot string `json:"before_snapshot"`
	AfterSnapshot  string `json:"after_snapshot"`
}

func (s *Server) createAuditRecord(w http.ResponseWriter, r *http.Request) {
	var req createAuditRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.CreateAuditRecord(model.AuditRecord{
		Operator:       req.Operator,
		Operation:      req.Operation,
		TargetType:     req.TargetType,
		TargetID:       req.TargetID,
		BeforeSnapshot: req.BeforeSnapshot,
		AfterSnapshot:  req.AfterSnapshot,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, item)
}

func (s *Server) listAuditRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AuditRecordFilter{
		Operator:   r.URL.Query().Get("operator"),
		Operation:  r.URL.Query().Get("operation"),
		TargetType: r.URL.Query().Get("target_type"),
		TargetID:   r.URL.Query().Get("target_id"),
		Keyword:    r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListAuditRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAuditRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.GetAuditRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) deleteAuditRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAuditRecord(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
