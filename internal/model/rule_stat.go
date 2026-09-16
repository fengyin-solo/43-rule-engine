package model

import (
	"time"
)

type RuleStat struct {
	ID             string    `json:"id"`
	RuleSetID      string    `json:"rule_set_id"`
	RuleID         string    `json:"rule_id"`
	Date           string    `json:"date"`
	TotalEval      int       `json:"total_eval"`
	HitCount       int       `json:"hit_count"`
	HitRatio       float64   `json:"hit_ratio"`
	AvgDurationMs  int       `json:"avg_duration_ms"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (rs *RuleStat) Validate() error {
	if rs.RuleSetID == "" {
		return NewValidationError("rule_set_id", "规则集 ID 不能为空")
	}
	if rs.RuleID == "" {
		return NewValidationError("rule_id", "规则 ID 不能为空")
	}
	if rs.Date == "" {
		return NewValidationError("date", "日期不能为空")
	}
	if rs.TotalEval < 0 {
		rs.TotalEval = 0
	}
	if rs.HitCount < 0 {
		rs.HitCount = 0
	}
	if rs.TotalEval > 0 {
		rs.HitRatio = float64(rs.HitCount) / float64(rs.TotalEval)
	} else {
		rs.HitRatio = 0
	}
	if rs.AvgDurationMs < 0 {
		rs.AvgDurationMs = 0
	}
	return nil
}

type RuleStatFilter struct {
	RuleSetID string
	RuleID    string
	Date      string
}

func (f RuleStatFilter) Match(rs *RuleStat) bool {
	if f.RuleSetID != "" && rs.RuleSetID != f.RuleSetID {
		return false
	}
	if f.RuleID != "" && rs.RuleID != f.RuleID {
		return false
	}
	if f.Date != "" && rs.Date != f.Date {
		return false
	}
	return true
}
