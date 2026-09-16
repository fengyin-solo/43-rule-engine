package handler

import (
	"net/http"

	"ruleengine/internal/model"
	"ruleengine/pkg/httpx"
)

func (s *Server) registerRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rules", s.createRule)
	mux.HandleFunc("GET /api/rules", s.listRules)
	mux.HandleFunc("GET /api/rules/{id}", s.getRule)
	mux.HandleFunc("PUT /api/rules/{id}", s.updateRule)
	mux.HandleFunc("DELETE /api/rules/{id}", s.deleteRule)
	mux.HandleFunc("POST /api/rules/{id}/toggle", s.toggleRule)
}

type createRuleRequest struct {
	RuleSetID    string   `json:"rule_set_id"`
	Name         string   `json:"name"`
	Priority     int      `json:"priority"`
	ConditionIDs []string `json:"condition_ids"`
	ActionIDs    []string `json:"action_ids"`
	Status       string   `json:"status"`
	Enabled      bool     `json:"enabled"`
}

func (s *Server) createRule(w http.ResponseWriter, r *http.Request) {
	var req createRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.CreateRule(model.Rule{
		RuleSetID:    req.RuleSetID,
		Name:         req.Name,
		Priority:     req.Priority,
		ConditionIDs: req.ConditionIDs,
		ActionIDs:    req.ActionIDs,
		Status:       req.Status,
		Enabled:      req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, item)
}

func (s *Server) listRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RuleFilter{
		RuleSetID: r.URL.Query().Get("rule_set_id"),
		Status:    r.URL.Query().Get("status"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	if v := r.URL.Query().Get("enabled"); v != "" {
		b := v == "true"
		filter.Enabled = &b
	}
	items, total, err := s.svc.ListRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.GetRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) updateRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.UpdateRule(id, model.Rule{
		RuleSetID:    req.RuleSetID,
		Name:         req.Name,
		Priority:     req.Priority,
		ConditionIDs: req.ConditionIDs,
		ActionIDs:    req.ActionIDs,
		Status:       req.Status,
		Enabled:      req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) deleteRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) toggleRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.ToggleRuleStatus(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}
