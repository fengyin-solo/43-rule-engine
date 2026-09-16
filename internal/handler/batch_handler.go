package handler

import (
	"net/http"
	"strconv"

	"ruleengine/pkg/httpx"
)

func (s *Server) registerBatchRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/batch/rules", s.batchCreateRules)
	mux.HandleFunc("POST /api/batch/conditions", s.batchCreateConditions)
	mux.HandleFunc("POST /api/batch/actions", s.batchCreateActions)
	mux.HandleFunc("POST /api/batch/rules/delete", s.batchDeleteRules)
}

type batchCreateRulesRequest struct {
	Rules []struct {
		RuleSetID    string   `json:"rule_set_id"`
		Name         string   `json:"name"`
		Priority     int      `json:"priority"`
		ConditionIDs []string `json:"condition_ids"`
		ActionIDs    []string `json:"action_ids"`
		Status       string   `json:"status"`
		Enabled      bool     `json:"enabled"`
	} `json:"rules"`
}

func (s *Server) batchCreateRules(w http.ResponseWriter, r *http.Request) {
	var req batchCreateRulesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	httpx.OK(w, map[string]interface{}{"count": len(req.Rules)})
}

type batchCreateConditionsRequest struct {
	Conditions []struct {
		Name        string `json:"name"`
		Field       string `json:"field"`
		Operator    string `json:"operator"`
		Value       string `json:"value"`
		ValueType   string `json:"value_type"`
		Conjunction string `json:"conjunction"`
	} `json:"conditions"`
}

func (s *Server) batchCreateConditions(w http.ResponseWriter, r *http.Request) {
	var req batchCreateConditionsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	httpx.OK(w, map[string]interface{}{"count": len(req.Conditions)})
}

type batchCreateActionsRequest struct {
	Actions []struct {
		Name       string            `json:"name"`
		ActionType string            `json:"action_type"`
		Params     map[string]string `json:"params"`
		Target     string            `json:"target"`
	} `json:"actions"`
}

func (s *Server) batchCreateActions(w http.ResponseWriter, r *http.Request) {
	var req batchCreateActionsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	httpx.OK(w, map[string]interface{}{"count": len(req.Actions)})
}

type batchDeleteRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchDeleteRules(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.BatchDeleteRules(req.IDs); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) registerReportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/reports/daily", s.getDailyReport)
	mux.HandleFunc("GET /api/reports/rule-sets/{id}", s.getRuleSetReport)
	mux.HandleFunc("GET /api/reports/trend", s.getExecutionTrend)
	mux.HandleFunc("GET /api/reports/condition-usage", s.getConditionUsage)
	mux.HandleFunc("GET /api/reports/action-usage", s.getActionUsage)
}

func (s *Server) getDailyReport(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = "2024-01-01"
	}
	data := s.svc.GenerateDailyReport(date)
	httpx.OK(w, data)
}

func (s *Server) getRuleSetReport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	if startDate == "" {
		startDate = "2024-01-01"
	}
	if endDate == "" {
		endDate = "2024-12-31"
	}
	data := s.svc.GenerateRuleSetReport(id, startDate, endDate)
	httpx.OK(w, data)
}

func (s *Server) getExecutionTrend(w http.ResponseWriter, r *http.Request) {
	ruleSetID := r.URL.Query().Get("rule_set_id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	data := s.svc.GetExecutionTrend(ruleSetID, limit)
	httpx.OK(w, data)
}

func (s *Server) getConditionUsage(w http.ResponseWriter, r *http.Request) {
	data := s.svc.GetConditionUsageReport()
	httpx.OK(w, data)
}

func (s *Server) getActionUsage(w http.ResponseWriter, r *http.Request) {
	data := s.svc.GetActionUsageReport()
	httpx.OK(w, data)
}
