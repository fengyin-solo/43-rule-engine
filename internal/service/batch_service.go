package service

import (
	"fmt"
	"time"

	"ruleengine/internal/model"
	"ruleengine/pkg/idgen"
)

func (s *Service) BatchCreateRules(inputs []model.Rule) ([]*model.Rule, error) {
	if len(inputs) == 0 {
		return nil, model.NewValidationError("inputs", "批量创建规则列表不能为空")
	}
	if len(inputs) > 100 {
		return nil, model.NewValidationError("inputs", "单次批量创建规则数量不能超过 100")
	}
	var results []*model.Rule
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, model.NewValidationError("inputs", fmt.Sprintf("规则 %s 校验失败: %v", input.Name, err))
		}
		if _, err := s.store.GetRuleSet(input.RuleSetID); err != nil {
			return nil, model.NewValidationError("rule_set_id", "所属规则集不存在: "+input.RuleSetID)
		}
		for _, cid := range input.ConditionIDs {
			if _, err := s.store.GetCondition(cid); err != nil {
				return nil, model.NewValidationError("condition_ids", "条件不存在: "+cid)
			}
		}
		for _, aid := range input.ActionIDs {
			if _, err := s.store.GetAction(aid); err != nil {
				return nil, model.NewValidationError("action_ids", "动作不存在: "+aid)
			}
		}
	}
	for _, input := range inputs {
		now := time.Now()
		r := &model.Rule{
			ID:           idgen.Hex(),
			RuleSetID:    input.RuleSetID,
			Name:         input.Name,
			Priority:     input.Priority,
			ConditionIDs: input.ConditionIDs,
			ActionIDs:    input.ActionIDs,
			Status:       input.Status,
			Enabled:      input.Enabled,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if err := s.store.CreateRule(r); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (s *Service) BatchCreateConditions(inputs []model.Condition) ([]*model.Condition, error) {
	if len(inputs) == 0 {
		return nil, model.NewValidationError("inputs", "批量创建条件列表不能为空")
	}
	if len(inputs) > 100 {
		return nil, model.NewValidationError("inputs", "单次批量创建条件数量不能超过 100")
	}
	var results []*model.Condition
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, model.NewValidationError("inputs", fmt.Sprintf("条件 %s 校验失败: %v", input.Name, err))
		}
	}
	for _, input := range inputs {
		now := time.Now()
		c := &model.Condition{
			ID:          idgen.Hex(),
			Name:        input.Name,
			Field:       input.Field,
			Operator:    input.Operator,
			Value:       input.Value,
			ValueType:   input.ValueType,
			Conjunction: input.Conjunction,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := s.store.CreateCondition(c); err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	return results, nil
}

func (s *Service) BatchCreateActions(inputs []model.Action) ([]*model.Action, error) {
	if len(inputs) == 0 {
		return nil, model.NewValidationError("inputs", "批量创建动作列表不能为空")
	}
	if len(inputs) > 100 {
		return nil, model.NewValidationError("inputs", "单次批量创建动作数量不能超过 100")
	}
	var results []*model.Action
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, model.NewValidationError("inputs", fmt.Sprintf("动作 %s 校验失败: %v", input.Name, err))
		}
	}
	for _, input := range inputs {
		now := time.Now()
		a := &model.Action{
			ID:         idgen.Hex(),
			Name:       input.Name,
			ActionType: input.ActionType,
			Params:     input.Params,
			Target:     input.Target,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err := s.store.CreateAction(a); err != nil {
			return nil, err
		}
		results = append(results, a)
	}
	return results, nil
}

func (s *Service) BatchDeleteRules(ids []string) error {
	if len(ids) == 0 {
		return model.NewValidationError("ids", "删除 ID 列表不能为空")
	}
	if len(ids) > 100 {
		return model.NewValidationError("ids", "单次批量删除数量不能超过 100")
	}
	for _, id := range ids {
		if err := s.store.DeleteRule(id); err != nil {
			return err
		}
	}
	return nil
}
