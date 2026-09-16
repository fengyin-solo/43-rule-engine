package service

import (
	"fmt"

	"ruleengine/internal/model"
)

func (s *Service) ValidateRuleSetExists(ruleSetID string) error {
	if _, err := s.store.GetRuleSet(ruleSetID); err != nil {
		return model.NewValidationError("rule_set_id", fmt.Sprintf("规则集 %s 不存在", ruleSetID))
	}
	return nil
}

func (s *Service) ValidateRuleExists(ruleID string) error {
	if _, err := s.store.GetRule(ruleID); err != nil {
		return model.NewValidationError("rule_id", fmt.Sprintf("规则 %s 不存在", ruleID))
	}
	return nil
}

func (s *Service) ValidateConditionExists(conditionID string) error {
	if _, err := s.store.GetCondition(conditionID); err != nil {
		return model.NewValidationError("condition_id", fmt.Sprintf("条件 %s 不存在", conditionID))
	}
	return nil
}

func (s *Service) ValidateActionExists(actionID string) error {
	if _, err := s.store.GetAction(actionID); err != nil {
		return model.NewValidationError("action_id", fmt.Sprintf("动作 %s 不存在", actionID))
	}
	return nil
}

func (s *Service) ValidateRuleSetIsPublished(ruleSetID string) error {
	rs, err := s.store.GetRuleSet(ruleSetID)
	if err != nil {
		return model.NewValidationError("rule_set_id", "规则集不存在")
	}
	if rs.Status != model.RuleSetStatusPublished {
		return model.NewValidationError("status", "规则集未发布")
	}
	return nil
}

func (s *Service) ValidateRuleSetCanModify(ruleSetID string) error {
	rs, err := s.store.GetRuleSet(ruleSetID)
	if err != nil {
		return model.NewValidationError("rule_set_id", "规则集不存在")
	}
	if rs.Status == model.RuleSetStatusDisabled {
		return model.NewValidationError("status", "规则集已禁用，无法修改")
	}
	if rs.Status == model.RuleSetStatusPublished {
		return model.NewValidationError("status", "规则集已发布，无法直接修改")
	}
	return nil
}

func (s *Service) ValidateConditionRefs(conditionIDs []string) error {
	for _, cid := range conditionIDs {
		if err := s.ValidateConditionExists(cid); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ValidateActionRefs(actionIDs []string) error {
	for _, aid := range actionIDs {
		if err := s.ValidateActionExists(aid); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ValidateRuleInput(input model.Rule) error {
	if err := input.Validate(); err != nil {
		return err
	}
	if err := s.ValidateRuleSetExists(input.RuleSetID); err != nil {
		return err
	}
	if err := s.ValidateConditionRefs(input.ConditionIDs); err != nil {
		return err
	}
	if err := s.ValidateActionRefs(input.ActionIDs); err != nil {
		return err
	}
	return nil
}
