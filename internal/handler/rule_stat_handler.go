package handler

import (
	"net/http"
	"strconv"

	"ruleengine/internal/model"
	"ruleengine/pkg/httpx"
)

func (s *Server) registerRuleStatRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rule-stats", s.createRuleStat)
	mux.HandleFunc("GET /api/rule-stats", s.listRuleStats)
	mux.HandleFunc("GET /api/rule-stats/{id}", s.getRuleStat)
	mux.HandleFunc("DELETE /api/rule-stats/{id}", s.deleteRuleStat)
}

type createRuleStatRequest struct {
	RuleSetID     string `json:"rule_set_id"`
	RuleID        string `json:"rule_id"`
	Date          string `json:"date"`
	TotalEval     int    `json:"total_eval"`
	HitCount      int    `json:"hit_count"`
	AvgDurationMs int    `json:"avg_duration_ms"`
}

func (s *Server) createRuleStat(w http.ResponseWriter, r *http.Request) {
	var req createRuleStatRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.CreateRuleStat(model.RuleStat{
		RuleSetID:     req.RuleSetID,
		RuleID:        req.RuleID,
		Date:          req.Date,
		TotalEval:     req.TotalEval,
		HitCount:      req.HitCount,
		AvgDurationMs: req.AvgDurationMs,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, item)
}

func (s *Server) listRuleStats(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RuleStatFilter{
		RuleSetID: r.URL.Query().Get("rule_set_id"),
		RuleID:    r.URL.Query().Get("rule_id"),
		Date:      r.URL.Query().Get("date"),
	}
	items, total, err := s.svc.ListRuleStats(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRuleStat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.svc.GetRuleStat(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

func (s *Server) deleteRuleStat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRuleStat(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.getStatsOverview)
	mux.HandleFunc("GET /api/stats/rule-sets/{id}", s.getRuleSetStats)
	mux.HandleFunc("GET /api/stats/top-rules", s.getTopRules)
}

func (s *Server) getStatsOverview(w http.ResponseWriter, r *http.Request) {
	data := s.svc.GetStatsOverview()
	httpx.OK(w, data)
}

func (s *Server) getRuleSetStats(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data := s.svc.GetRuleSetStats(id)
	httpx.OK(w, data)
}

func (s *Server) getTopRules(w http.ResponseWriter, r *http.Request) {
	n, _ := strconv.Atoi(r.URL.Query().Get("n"))
	if n < 1 {
		n = 5
	}
	items := s.svc.GetTopNRulesByHit(n)
	httpx.OK(w, items)
}
