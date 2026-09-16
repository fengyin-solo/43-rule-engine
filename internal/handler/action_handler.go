package handler

import (
	"net/http"

	"ruleengine/internal/model"
	"ruleengine/pkg/httpx"
)

func (s *Server) registerActionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/actions", s.createAction)
	mux.HandleFunc("GET /api/actions", s.listActions)
	mux.HandleFunc("GET /api/actions/{id}", s.getAction)
	mux.HandleFunc("PUT /api/actions/{id}", s.updateAction)
	mux.HandleFunc("DELETE /api/actions/{id}", s.deleteAction)
}

type createActionRequest struct {
	Name       string            `json:"name"`
	ActionType string            `json:"action_type"`
	Params     map[string]string `json:"params"`
	Target     string            `json:"target"`
}

func (s *Server) createAction(w http.ResponseWriter, r *http.Request) {
	var req createActionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.CreateAction(model.Action{
		Name:       req.Name,
		ActionType: req.ActionType,
		Params:     req.Params,
		Target:     req.Target,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, item)
}

func (s *Server) listActions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ActionFilter{
		ActionType: r.URL.Query().Get("action_type"),
		Keyword:    r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListActions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.GetAction(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) updateAction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createActionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.UpdateAction(id, model.Action{
		Name:       req.Name,
		ActionType: req.ActionType,
		Params:     req.Params,
		Target:     req.Target,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) deleteAction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAction(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
