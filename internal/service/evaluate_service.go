package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"ruleengine/internal/model"
	"ruleengine/pkg/idgen"
)

func (s *Service) EvaluateRuleSet(ruleSetID string, input map[string]interface{}) (*model.EvaluationRecord, error) {
	rs, err := s.store.GetRuleSet(ruleSetID)
	if err != nil {
		return nil, err
	}
	if rs.Status != model.RuleSetStatusPublished {
		return nil, model.NewValidationError("status", "只有已发布的规则集才能评估")
	}
	allRules := s.store.ListRulesByRuleSetID(ruleSetID)
	var activeRules []*model.Rule
	for _, r := range allRules {
		if r.Status == model.RuleStatusActive && r.Enabled {
			activeRules = append(activeRules, r)
		}
	}
	startAll := time.Now()
	var matchedRuleIDs []string
	var results []map[string]interface{}
	for _, r := range activeRules {
		startRule := time.Now()
		hit, err := s.evaluateRule(r, input)
		duration := int(time.Since(startRule).Milliseconds())
		if err != nil {
			s.log.Warnf("评估规则 %s 出错: %v", r.ID, err)
			continue
		}
		el := &model.ExecutionLog{
			ID:         idgen.Hex(),
			RuleSetID:  ruleSetID,
			RuleID:     r.ID,
			InputKey:   fmt.Sprintf("%v", input),
			DurationMs: duration,
			Hit:        hit,
			ExecutedAt: time.Now(),
		}
		_ = s.store.CreateExecutionLog(el)
		if hit {
			matchedRuleIDs = append(matchedRuleIDs, r.ID)
			actionResults := s.executeActions(r.ActionIDs)
			results = append(results, map[string]interface{}{
				"rule_id": r.ID,
				"actions": actionResults,
			})
		}
	}
	totalDuration := int(time.Since(startAll).Milliseconds())
	result := model.ResultMiss
	if len(matchedRuleIDs) > 0 {
		result = model.ResultHit
	}
	er := &model.EvaluationRecord{
		ID:             idgen.Hex(),
		RuleSetID:      ruleSetID,
		InputSnapshot:  input,
		MatchedRuleIDs: matchedRuleIDs,
		Result:         result,
		EvaluatedAt:    time.Now(),
	}
	if err := s.store.CreateEvaluationRecord(er); err != nil {
		return nil, err
	}
	if len(matchedRuleIDs) > 0 {
		s.log.Infof("规则集 %s 评估命中 %d 条规则", ruleSetID, len(matchedRuleIDs))
	} else {
		s.log.Infof("规则集 %s 评估未命中", ruleSetID)
	}
	_ = totalDuration
	_ = results
	return er, nil
}

func (s *Service) evaluateRule(r *model.Rule, input map[string]interface{}) (bool, error) {
	if len(r.ConditionIDs) == 0 {
		return true, nil
	}
	var results []bool
	var conjunction string
	for _, cid := range r.ConditionIDs {
		c, err := s.store.GetCondition(cid)
		if err != nil {
			return false, err
		}
		if conjunction == "" {
			conjunction = c.Conjunction
		}
		val, ok := input[c.Field]
		if !ok {
			results = append(results, false)
			continue
		}
		matched := evaluateCondition(c, val)
		results = append(results, matched)
	}
	if len(results) == 0 {
		return true, nil
	}
	return combineResults(results, conjunction), nil
}

func evaluateCondition(c *model.Condition, inputValue interface{}) bool {
	switch c.Operator {
	case model.OpEQ:
		return compareEqual(c.ValueType, c.Value, inputValue)
	case model.OpNE:
		return !compareEqual(c.ValueType, c.Value, inputValue)
	case model.OpGT:
		return compareNumber(c.Value, inputValue, func(a, b float64) bool { return a > b })
	case model.OpGE:
		return compareNumber(c.Value, inputValue, func(a, b float64) bool { return a >= b })
	case model.OpLT:
		return compareNumber(c.Value, inputValue, func(a, b float64) bool { return a < b })
	case model.OpLE:
		return compareNumber(c.Value, inputValue, func(a, b float64) bool { return a <= b })
	case model.OpContains:
		return compareContains(c.ValueType, c.Value, inputValue)
	case model.OpIn:
		return compareIn(c.ValueType, c.Value, inputValue)
	case model.OpBetween:
		return compareBetween(c.ValueType, c.Value, inputValue)
	}
	return false
}

func compareEqual(valueType, condValue string, inputValue interface{}) bool {
	switch valueType {
	case model.ValueTypeNumber:
		a, err1 := strconv.ParseFloat(condValue, 64)
		b, err2 := toFloat64(inputValue)
		if err1 != nil || err2 != nil {
			return false
		}
		return a == b
	case model.ValueTypeBool:
		a, err1 := strconv.ParseBool(condValue)
		b, err2 := toBool(inputValue)
		if err1 != nil || err2 != nil {
			return false
		}
		return a == b
	default:
		return fmt.Sprintf("%v", inputValue) == condValue
	}
}

func compareNumber(condValue string, inputValue interface{}, cmp func(a, b float64) bool) bool {
	a, err1 := strconv.ParseFloat(condValue, 64)
	b, err2 := toFloat64(inputValue)
	if err1 != nil || err2 != nil {
		return false
	}
	return cmp(b, a)
}

func compareContains(valueType, condValue string, inputValue interface{}) bool {
	str := fmt.Sprintf("%v", inputValue)
	return strings.Contains(str, condValue)
}

func compareIn(valueType, condValue string, inputValue interface{}) bool {
	parts := strings.Split(condValue, ",")
	inputStr := fmt.Sprintf("%v", inputValue)
	for _, p := range parts {
		if strings.TrimSpace(p) == inputStr {
			return true
		}
	}
	return false
}

func compareBetween(valueType, condValue string, inputValue interface{}) bool {
	parts := strings.Split(condValue, ",")
	if len(parts) != 2 {
		return false
	}
	minVal, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	maxVal, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	b, err3 := toFloat64(inputValue)
	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}
	return b >= minVal && b <= maxVal
}

func combineResults(results []bool, conjunction string) bool {
	if len(results) == 0 {
		return true
	}
	if conjunction == model.ConjunctionOr {
		for _, r := range results {
			if r {
				return true
			}
		}
		return false
	}
	for _, r := range results {
		if !r {
			return false
		}
	}
	return true
}

func toFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	case json.Number:
		return val.Float64()
	default:
		return strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
	}
}

func toBool(v interface{}) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case string:
		return strconv.ParseBool(val)
	default:
		return strconv.ParseBool(fmt.Sprintf("%v", v))
	}
}

func (s *Service) executeActions(actionIDs []string) []map[string]interface{} {
	var results []map[string]interface{}
	for _, aid := range actionIDs {
		a, err := s.store.GetAction(aid)
		if err != nil {
			continue
		}
		results = append(results, map[string]interface{}{
			"action_id":   a.ID,
			"action_type": a.ActionType,
			"target":      a.Target,
			"params":      a.Params,
		})
	}
	return results
}
