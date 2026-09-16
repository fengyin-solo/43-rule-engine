package model

import (
	"strings"
	"time"
)

type ExecutionLog struct {
	ID          string    `json:"id"`
	RuleSetID   string    `json:"rule_set_id"`
	RuleID      string    `json:"rule_id"`
	InputKey    string    `json:"input_key"`
	DurationMs  int       `json:"duration_ms"`
	Hit         bool      `json:"hit"`
	ExecutedAt  time.Time `json:"executed_at"`
}

func (el *ExecutionLog) Validate() error {
	if el.RuleSetID == "" {
		return NewValidationError("rule_set_id", "规则集 ID 不能为空")
	}
	if el.RuleID == "" {
		return NewValidationError("rule_id", "规则 ID 不能为空")
	}
	if el.DurationMs < 0 {
		el.DurationMs = 0
	}
	return nil
}

type ExecutionLogFilter struct {
	RuleSetID string
	RuleID    string
	Hit       *bool
	Keyword   string
}

func (f ExecutionLogFilter) Match(el *ExecutionLog) bool {
	if f.RuleSetID != "" && el.RuleSetID != f.RuleSetID {
		return false
	}
	if f.RuleID != "" && el.RuleID != f.RuleID {
		return false
	}
	if f.Hit != nil && el.Hit != *f.Hit {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(el.InputKey), k) {
			return false
		}
	}
	return true
}
