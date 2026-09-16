package handler

import (
	"net/http"

	"ruleengine/pkg/httpx"
)

func (s *Server) registerExportImportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/rule-sets/{id}/export", s.exportRuleSet)
	mux.HandleFunc("POST /api/rule-sets/import", s.importRuleSet)
}

func (s *Server) exportRuleSet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := s.svc.ExportRuleSetSnapshot(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=ruleset-"+id+".json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(data))
}

type importRequest struct {
	Data     string `json:"data"`
	Operator string `json:"operator"`
}

func (s *Server) importRuleSet(w http.ResponseWriter, r *http.Request) {
	var req importRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if req.Operator == "" {
		req.Operator = "system"
	}
	rs, err := s.svc.ImportRuleSetSnapshot(req.Data, req.Operator)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rs)
}
