package model

import (
	"strings"
	"time"
)

const (
	RuleVersionStatusActive   = "active"
	RuleVersionStatusArchived = "archived"
	RuleVersionStatusRolledBack = "rolled_back"
)

var ruleVersionTransitions = map[string]map[string]bool{
	RuleVersionStatusActive:     {RuleVersionStatusArchived: true, RuleVersionStatusRolledBack: true},
	RuleVersionStatusArchived:   {},
	RuleVersionStatusRolledBack: {},
}

func RuleVersionCanTransition(from, to string) bool {
	if m, ok := ruleVersionTransitions[from]; ok {
		return m[to]
	}
	return false
}

type RuleVersion struct {
	ID          string    `json:"id"`
	RuleSetID   string    `json:"rule_set_id"`
	Version     int       `json:"version"`
	ChangeNote  string    `json:"change_note"`
	ChangedBy   string    `json:"changed_by"`
	PublishedAt time.Time `json:"published_at"`
	Status      string    `json:"status"`
	Snapshot    string    `json:"snapshot"`
	CreatedAt   time.Time `json:"created_at"`
}

func (rv *RuleVersion) Validate() error {
	rv.ChangeNote = strings.TrimSpace(rv.ChangeNote)
	rv.ChangedBy = strings.TrimSpace(rv.ChangedBy)
	if rv.RuleSetID == "" {
		return NewValidationError("rule_set_id", "规则集 ID 不能为空")
	}
	if rv.Version < 1 {
		return NewValidationError("version", "版本号必须大于 0")
	}
	if rv.Status == "" {
		rv.Status = RuleVersionStatusActive
	}
	if rv.Status != RuleVersionStatusActive && rv.Status != RuleVersionStatusArchived && rv.Status != RuleVersionStatusRolledBack {
		return NewValidationError("status", "版本状态不合法")
	}
	return nil
}

type RuleVersionFilter struct {
	RuleSetID string
	Status    string
	Keyword   string
}

func (f RuleVersionFilter) Match(rv *RuleVersion) bool {
	if f.RuleSetID != "" && rv.RuleSetID != f.RuleSetID {
		return false
	}
	if f.Status != "" && rv.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(rv.ChangeNote), k) {
			return false
		}
	}
	return true
}
