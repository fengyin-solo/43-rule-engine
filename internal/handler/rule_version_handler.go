package handler

import (
	"net/http"

	"ruleengine/internal/model"
	"ruleengine/pkg/httpx"
)

func (s *Server) registerRuleVersionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rule-versions", s.createRuleVersion)
	mux.HandleFunc("GET /api/rule-versions", s.listRuleVersions)
	mux.HandleFunc("GET /api/rule-versions/{id}", s.getRuleVersion)
	mux.HandleFunc("PUT /api/rule-versions/{id}", s.updateRuleVersion)
	mux.HandleFunc("DELETE /api/rule-versions/{id}", s.deleteRuleVersion)
	mux.HandleFunc("POST /api/rule-versions/{id}/rollback", s.rollbackRuleVersion)
}

type createRuleVersionRequest struct {
	RuleSetID  string `json:"rule_set_id"`
	Version    int    `json:"version"`
	ChangeNote string `json:"change_note"`
	ChangedBy  string `json:"changed_by"`
	Status     string `json:"status"`
	Snapshot   string `json:"snapshot"`
}

func (s *Server) createRuleVersion(w http.ResponseWriter, r *http.Request) {
	var req createRuleVersionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.CreateRuleVersion(model.RuleVersion{
		RuleSetID:  req.RuleSetID,
		Version:    req.Version,
		ChangeNote: req.ChangeNote,
		ChangedBy:  req.ChangedBy,
		Status:     req.Status,
		Snapshot:   req.Snapshot,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, item)
}

func (s *Server) listRuleVersions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RuleVersionFilter{
		RuleSetID: r.URL.Query().Get("rule_set_id"),
		Status:    r.URL.Query().Get("status"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRuleVersions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRuleVersion(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.GetRuleVersion(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) updateRuleVersion(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createRuleVersionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.UpdateRuleVersion(id, model.RuleVersion{
		ChangeNote: req.ChangeNote,
		ChangedBy:  req.ChangedBy,
		Status:     req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) deleteRuleVersion(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRuleVersion(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) rollbackRuleVersion(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	operator := r.Header.Get("X-Operator")
	if operator == "" {
		operator = "system"
	}
	item, err := s.svc.RollbackRuleVersion(id, operator)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}
