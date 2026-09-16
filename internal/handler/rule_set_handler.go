package handler

import (
	"net/http"

	"ruleengine/internal/model"
	"ruleengine/pkg/httpx"
)

func (s *Server) registerRuleSetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rule-sets", s.createRuleSet)
	mux.HandleFunc("GET /api/rule-sets", s.listRuleSets)
	mux.HandleFunc("GET /api/rule-sets/{id}", s.getRuleSet)
	mux.HandleFunc("PUT /api/rule-sets/{id}", s.updateRuleSet)
	mux.HandleFunc("DELETE /api/rule-sets/{id}", s.deleteRuleSet)
	mux.HandleFunc("POST /api/rule-sets/{id}/publish", s.publishRuleSet)
	mux.HandleFunc("POST /api/rule-sets/{id}/disable", s.disableRuleSet)
}

type createRuleSetRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) createRuleSet(w http.ResponseWriter, r *http.Request) {
	var req createRuleSetRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rs, err := s.svc.CreateRuleSet(model.RuleSet{Name: req.Name, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rs)
}

func (s *Server) listRuleSets(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RuleSetFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRuleSets(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRuleSet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rs, err := s.svc.GetRuleSet(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rs)
}

func (s *Server) updateRuleSet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createRuleSetRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rs, err := s.svc.UpdateRuleSet(id, model.RuleSet{Name: req.Name, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rs)
}

func (s *Server) deleteRuleSet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRuleSet(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) publishRuleSet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	changedBy := r.Header.Get("X-Operator")
	if changedBy == "" {
		changedBy = "system"
	}
	rs, err := s.svc.PublishRuleSet(id, changedBy)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rs)
}

func (s *Server) disableRuleSet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	changedBy := r.Header.Get("X-Operator")
	if changedBy == "" {
		changedBy = "system"
	}
	rs, err := s.svc.DisableRuleSet(id, changedBy)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rs)
}
