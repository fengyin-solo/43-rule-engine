package model

import (
	"strings"
	"time"
)

const (
	ResultHit  = "hit"
	ResultMiss = "miss"
)

type EvaluationRecord struct {
	ID            string                 `json:"id"`
	RuleSetID     string                 `json:"rule_set_id"`
	InputSnapshot map[string]interface{} `json:"input_snapshot"`
	MatchedRuleIDs []string              `json:"matched_rule_ids"`
	Result        string                 `json:"result"`
	EvaluatedAt   time.Time              `json:"evaluated_at"`
}

func (er *EvaluationRecord) Validate() error {
	if er.RuleSetID == "" {
		return NewValidationError("rule_set_id", "规则集 ID 不能为空")
	}
	if er.Result == "" {
		er.Result = ResultMiss
	}
	if er.Result != ResultHit && er.Result != ResultMiss {
		return NewValidationError("result", "评估结果不合法")
	}
	return nil
}

type EvaluationRecordFilter struct {
	RuleSetID string
	Result    string
	Keyword   string
}

func (f EvaluationRecordFilter) Match(er *EvaluationRecord) bool {
	if f.RuleSetID != "" && er.RuleSetID != f.RuleSetID {
		return false
	}
	if f.Result != "" && er.Result != f.Result {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" {
			found := false
			for key := range er.InputSnapshot {
				if strings.Contains(strings.ToLower(key), k) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}
	return true
}
