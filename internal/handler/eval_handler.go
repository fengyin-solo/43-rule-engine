package handler

import (
	"net/http"

	"ruleengine/pkg/httpx"
)

func (s *Server) registerEvalRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/evaluate", s.evaluate)
}

type evaluateRequest struct {
	RuleSetID string                 `json:"rule_set_id"`
	Input     map[string]interface{} `json:"input"`
}

func (s *Server) evaluate(w http.ResponseWriter, r *http.Request) {
	var req evaluateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if req.RuleSetID == "" {
		httpx.BadRequest(w, "rule_set_id 不能为空")
		return
	}
	result, err := s.svc.EvaluateRuleSet(req.RuleSetID, req.Input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
