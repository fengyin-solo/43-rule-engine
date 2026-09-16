package model

import (
	"strings"
	"time"
)

const (
	RuleStatusActive   = "active"
	RuleStatusInactive = "inactive"
)

type Rule struct {
	ID          string    `json:"id"`
	RuleSetID   string    `json:"rule_set_id"`
	Name        string    `json:"name"`
	Priority    int       `json:"priority"`
	ConditionIDs []string `json:"condition_ids"`
	ActionIDs   []string  `json:"action_ids"`
	Status      string    `json:"status"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r *Rule) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return NewValidationError("name", "规则名称不能为空")
	}
	if r.RuleSetID == "" {
		return NewValidationError("rule_set_id", "规则所属规则集不能为空")
	}
	if r.Status == "" {
		r.Status = RuleStatusActive
	}
	if r.Status != RuleStatusActive && r.Status != RuleStatusInactive {
		return NewValidationError("status", "规则状态不合法")
	}
	return nil
}

type RuleFilter struct {
	RuleSetID string
	Name      string
	Status    string
	Enabled   *bool
	Keyword   string
}

func (f RuleFilter) Match(r *Rule) bool {
	if f.RuleSetID != "" && r.RuleSetID != f.RuleSetID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Enabled != nil && r.Enabled != *f.Enabled {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Name), k) {
			return false
		}
	}
	return true
}
