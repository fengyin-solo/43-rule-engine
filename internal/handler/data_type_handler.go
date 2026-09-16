package handler

import (
	"encoding/json"
	"net/http"

	"ruleengine/internal/model"
	"ruleengine/pkg/httpx"
)

func (s *Server) registerDataTypeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/data-types", s.createDataType)
	mux.HandleFunc("GET /api/data-types", s.listDataTypes)
	mux.HandleFunc("GET /api/data-types/{id}", s.getDataType)
	mux.HandleFunc("PUT /api/data-types/{id}", s.updateDataType)
	mux.HandleFunc("DELETE /api/data-types/{id}", s.deleteDataType)
}

type createDataTypeRequest struct {
	Name        string          `json:"name"`
	Fields      json.RawMessage `json:"fields"`
	Description string          `json:"description"`
	Enabled     bool            `json:"enabled"`
}

func (s *Server) createDataType(w http.ResponseWriter, r *http.Request) {
	var req createDataTypeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.CreateDataType(model.DataType{
		Name:        req.Name,
		Fields:      req.Fields,
		Description: req.Description,
		Enabled:     req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, item)
}

func (s *Server) listDataTypes(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DataTypeFilter{
		Keyword: r.URL.Query().Get("keyword"),
	}
	if v := r.URL.Query().Get("enabled"); v != "" {
		b := v == "true"
		filter.Enabled = &b
	}
	items, total, err := s.svc.ListDataTypes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDataType(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.GetDataType(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) updateDataType(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createDataTypeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.UpdateDataType(id, model.DataType{
		Name:        req.Name,
		Fields:      req.Fields,
		Description: req.Description,
		Enabled:     req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) deleteDataType(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteDataType(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
