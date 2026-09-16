package service

import (
	"fmt"
	"sort"
	"time"

	"ruleengine/internal/model"
	"ruleengine/pkg/idgen"
)

func (s *Service) CreateRuleStat(input model.RuleStat) (*model.RuleStat, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	rs := &model.RuleStat{
		ID:            idgen.Hex(),
		RuleSetID:     input.RuleSetID,
		RuleID:        input.RuleID,
		Date:          input.Date,
		TotalEval:     input.TotalEval,
		HitCount:      input.HitCount,
		HitRatio:      input.HitRatio,
		AvgDurationMs: input.AvgDurationMs,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.store.CreateRuleStat(rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (s *Service) GetRuleStat(id string) (*model.RuleStat, error) {
	return s.store.GetRuleStat(id)
}

func (s *Service) ListRuleStats(filter model.RuleStatFilter, page, size int) ([]*model.RuleStat, int, error) {
	all := s.store.ListRuleStats()
	matched := make([]*model.RuleStat, 0, len(all))
	for _, rs := range all {
		if filter.Match(rs) {
			matched = append(matched, rs)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Date > matched[j].Date
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RuleStat{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteRuleStat(id string) error {
	return s.store.DeleteRuleStat(id)
}

func (s *Service) GetStatsOverview() map[string]interface{} {
	allRuleSets := s.store.ListRuleSets()
	publishedCount := 0
	for _, rs := range allRuleSets {
		if rs.Status == model.RuleSetStatusPublished {
			publishedCount++
		}
	}
	allEvals := s.store.ListEvaluationRecords()
	totalEval := len(allEvals)
	totalHit := 0
	for _, er := range allEvals {
		if er.Result == model.ResultHit {
			totalHit++
		}
	}
	avgHitRatio := 0.0
	if totalEval > 0 {
		avgHitRatio = float64(totalHit) / float64(totalEval)
	}
	allLogs := s.store.ListExecutionLogs()
	totalDuration := 0
	for _, el := range allLogs {
		totalDuration += el.DurationMs
	}
	avgDuration := 0
	if len(allLogs) > 0 {
		avgDuration = totalDuration / len(allLogs)
	}
	return map[string]interface{}{
		"rule_set_count":    len(allRuleSets),
		"published_count":   publishedCount,
		"total_evaluation":  totalEval,
		"avg_hit_ratio":     fmt.Sprintf("%.2f", avgHitRatio),
		"avg_duration_ms":   avgDuration,
	}
}

func (s *Service) GetRuleSetStats(ruleSetID string) map[string]interface{} {
	all := s.store.ListRuleStats()
	var matched []*model.RuleStat
	for _, rs := range all {
		if rs.RuleSetID == ruleSetID {
			matched = append(matched, rs)
		}
	}
	totalEval := 0
	totalHit := 0
	totalDuration := 0
	for _, rs := range matched {
		totalEval += rs.TotalEval
		totalHit += rs.HitCount
		totalDuration += rs.AvgDurationMs * rs.TotalEval
	}
	avgDuration := 0
	if totalEval > 0 {
		avgDuration = totalDuration / totalEval
	}
	return map[string]interface{}{
		"rule_set_id":      ruleSetID,
		"total_eval":       totalEval,
		"total_hit":        totalHit,
		"avg_duration_ms":  avgDuration,
		"daily_stats":      matched,
	}
}

func (s *Service) GetTopNRulesByHit(n int) []*model.RuleStat {
	all := s.store.ListRuleStats()
	group := make(map[string]*model.RuleStat)
	for _, rs := range all {
		key := rs.RuleID
		if exist, ok := group[key]; ok {
			exist.TotalEval += rs.TotalEval
			exist.HitCount += rs.HitCount
		} else {
			cp := *rs
			group[key] = &cp
		}
	}
	list := make([]*model.RuleStat, 0, len(group))
	for _, rs := range group {
		if rs.TotalEval > 0 {
			rs.HitRatio = float64(rs.HitCount) / float64(rs.TotalEval)
		}
		list = append(list, rs)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].HitCount > list[j].HitCount
	})
	if n > len(list) {
		n = len(list)
	}
	return list[:n]
}
