package service

import (
	"sort"
	"time"

	"ruleengine/internal/model"
	"ruleengine/pkg/idgen"
)

func (s *Service) CreateRule(input model.Rule) (*model.Rule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetRuleSet(input.RuleSetID); err != nil {
		return nil, model.NewValidationError("rule_set_id", "所属规则集不存在")
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
	return r, nil
}

func (s *Service) GetRule(id string) (*model.Rule, error) {
	return s.store.GetRule(id)
}

func (s *Service) ListRules(filter model.RuleFilter, page, size int) ([]*model.Rule, int, error) {
	all := s.store.ListRules()
	matched := make([]*model.Rule, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].Priority != matched[j].Priority {
			return matched[i].Priority > matched[j].Priority
		}
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Rule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRule(id string, input model.Rule) (*model.Rule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	r, err := s.store.GetRule(id)
	if err != nil {
		return nil, err
	}
	if input.RuleSetID != "" && input.RuleSetID != r.RuleSetID {
		if _, err := s.store.GetRuleSet(input.RuleSetID); err != nil {
			return nil, model.NewValidationError("rule_set_id", "所属规则集不存在")
		}
		r.RuleSetID = input.RuleSetID
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
	r.Name = input.Name
	r.Priority = input.Priority
	r.ConditionIDs = input.ConditionIDs
	r.ActionIDs = input.ActionIDs
	r.Status = input.Status
	r.Enabled = input.Enabled
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateRule(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteRule(id string) error {
	return s.store.DeleteRule(id)
}

func (s *Service) ToggleRuleStatus(id string) (*model.Rule, error) {
	r, err := s.store.GetRule(id)
	if err != nil {
		return nil, err
	}
	if r.Status == model.RuleStatusActive {
		r.Status = model.RuleStatusInactive
	} else {
		r.Status = model.RuleStatusActive
	}
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateRule(r); err != nil {
		return nil, err
	}
	return r, nil
}
