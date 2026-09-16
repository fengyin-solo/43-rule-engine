package service

import (
	"fmt"
	"sort"
	"time"

	"ruleengine/internal/model"
)

func (s *Service) GenerateDailyReport(date string) map[string]interface{} {
	allStats := s.store.ListRuleStats()
	var dailyStats []*model.RuleStat
	for _, st := range allStats {
		if st.Date == date {
			dailyStats = append(dailyStats, st)
		}
	}
	totalEval := 0
	totalHit := 0
	totalDuration := 0
	for _, st := range dailyStats {
		totalEval += st.TotalEval
		totalHit += st.HitCount
		totalDuration += st.AvgDurationMs * st.TotalEval
	}
	avgDuration := 0
	if totalEval > 0 {
		avgDuration = totalDuration / totalEval
	}
	hitRatio := 0.0
	if totalEval > 0 {
		hitRatio = float64(totalHit) / float64(totalEval)
	}
	return map[string]interface{}{
		"date":            date,
		"total_eval":      totalEval,
		"total_hit":       totalHit,
		"hit_ratio":       fmt.Sprintf("%.4f", hitRatio),
		"avg_duration_ms": avgDuration,
		"rule_count":      len(dailyStats),
	}
}

func (s *Service) GenerateRuleSetReport(ruleSetID string, startDate, endDate string) map[string]interface{} {
	allStats := s.store.ListRuleStats()
	var matched []*model.RuleStat
	for _, st := range allStats {
		if st.RuleSetID == ruleSetID && st.Date >= startDate && st.Date <= endDate {
			matched = append(matched, st)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Date < matched[j].Date
	})
	totalEval := 0
	totalHit := 0
	for _, st := range matched {
		totalEval += st.TotalEval
		totalHit += st.HitCount
	}
	return map[string]interface{}{
		"rule_set_id": ruleSetID,
		"start_date":  startDate,
		"end_date":    endDate,
		"total_eval":  totalEval,
		"total_hit":   totalHit,
		"items":       matched,
	}
}

func (s *Service) GetExecutionTrend(ruleSetID string, limit int) []map[string]interface{} {
	if limit < 1 {
		limit = 7
	}
	if limit > 30 {
		limit = 30
	}
	allLogs := s.store.ListExecutionLogs()
	var matched []*model.ExecutionLog
	for _, el := range allLogs {
		if el.RuleSetID == ruleSetID {
			matched = append(matched, el)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].ExecutedAt.After(matched[j].ExecutedAt)
	})
	if len(matched) > limit {
		matched = matched[:limit]
	}
	var results []map[string]interface{}
	for _, el := range matched {
		results = append(results, map[string]interface{}{
			"rule_id":     el.RuleID,
			"duration_ms": el.DurationMs,
			"hit":         el.Hit,
			"executed_at": el.ExecutedAt.Format(time.RFC3339),
		})
	}
	return results
}

func (s *Service) GetConditionUsageReport() []map[string]interface{} {
	allRules := s.store.ListRules()
	usage := make(map[string]int)
	for _, r := range allRules {
		for _, cid := range r.ConditionIDs {
			usage[cid]++
		}
	}
	type pair struct {
		id    string
		count int
	}
	var pairs []pair
	for id, count := range usage {
		pairs = append(pairs, pair{id: id, count: count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})
	var results []map[string]interface{}
	for _, p := range pairs {
		c, err := s.store.GetCondition(p.id)
		name := p.id
		if err == nil {
			name = c.Name
		}
		results = append(results, map[string]interface{}{
			"condition_id": p.id,
			"name":         name,
			"used_count":   p.count,
		})
	}
	return results
}

func (s *Service) GetActionUsageReport() []map[string]interface{} {
	allRules := s.store.ListRules()
	usage := make(map[string]int)
	for _, r := range allRules {
		for _, aid := range r.ActionIDs {
			usage[aid]++
		}
	}
	type pair struct {
		id    string
		count int
	}
	var pairs []pair
	for id, count := range usage {
		pairs = append(pairs, pair{id: id, count: count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})
	var results []map[string]interface{}
	for _, p := range pairs {
		a, err := s.store.GetAction(p.id)
		name := p.id
		if err == nil {
			name = a.Name
		}
		results = append(results, map[string]interface{}{
			"action_id":  p.id,
			"name":       name,
			"used_count": p.count,
		})
	}
	return results
}
