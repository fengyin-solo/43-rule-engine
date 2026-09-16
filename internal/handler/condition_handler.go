package handler

import (
	"net/http"

	"ruleengine/internal/model"
	"ruleengine/pkg/httpx"
)

func (s *Server) registerConditionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/conditions", s.createCondition)
	mux.HandleFunc("GET /api/conditions", s.listConditions)
	mux.HandleFunc("GET /api/conditions/{id}", s.getCondition)
	mux.HandleFunc("PUT /api/conditions/{id}", s.updateCondition)
	mux.HandleFunc("DELETE /api/conditions/{id}", s.deleteCondition)
}

type createConditionRequest struct {
	Name        string `json:"name"`
	Field       string `json:"field"`
	Operator    string `json:"operator"`
	Value       string `json:"value"`
	ValueType   string `json:"value_type"`
	Conjunction string `json:"conjunction"`
}

func (s *Server) createCondition(w http.ResponseWriter, r *http.Request) {
	var req createConditionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.CreateCondition(model.Condition{
		Name:        req.Name,
		Field:       req.Field,
		Operator:    req.Operator,
		Value:       req.Value,
		ValueType:   req.ValueType,
		Conjunction: req.Conjunction,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, item)
}

func (s *Server) listConditions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ConditionFilter{
		Operator:  r.URL.Query().Get("operator"),
		ValueType: r.URL.Query().Get("value_type"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListConditions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCondition(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.GetCondition(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) updateCondition(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createConditionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.UpdateCondition(id, model.Condition{
		Name:        req.Name,
		Field:       req.Field,
		Operator:    req.Operator,
		Value:       req.Value,
		ValueType:   req.ValueType,
		Conjunction: req.Conjunction,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) deleteCondition(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCondition(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
