package model

import (
	"strings"
	"time"
)

const (
	RuleSetStatusDraft     = "draft"
	RuleSetStatusPublished = "published"
	RuleSetStatusDisabled  = "disabled"
)

var ruleSetTransitions = map[string]map[string]bool{
	RuleSetStatusDraft:     {RuleSetStatusPublished: true},
	RuleSetStatusPublished: {RuleSetStatusDisabled: true},
	RuleSetStatusDisabled:  {},
}

func RuleSetCanTransition(from, to string) bool {
	if m, ok := ruleSetTransitions[from]; ok {
		return m[to]
	}
	return false
}

type RuleSet struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     int       `json:"version"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (rs *RuleSet) Validate() error {
	rs.Name = strings.TrimSpace(rs.Name)
	rs.Description = strings.TrimSpace(rs.Description)
	if rs.Name == "" {
		return NewValidationError("name", "规则集名称不能为空")
	}
	if len(rs.Name) > 128 {
		return NewValidationError("name", "规则集名称长度不能超过 128")
	}
	if rs.Status == "" {
		rs.Status = RuleSetStatusDraft
	}
	if rs.Status != RuleSetStatusDraft && rs.Status != RuleSetStatusPublished && rs.Status != RuleSetStatusDisabled {
		return NewValidationError("status", "规则集状态不合法")
	}
	if rs.Version < 1 {
		rs.Version = 1
	}
	return nil
}

type RuleSetFilter struct {
	Name    string
	Status  string
	Keyword string
}

func (f RuleSetFilter) Match(rs *RuleSet) bool {
	if f.Status != "" && rs.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(rs.Name), k) &&
			!strings.Contains(strings.ToLower(rs.Description), k) {
			return false
		}
	}
	return true
}
