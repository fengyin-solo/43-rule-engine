package service

import (
	"sort"
	"time"

	"ruleengine/internal/model"
	"ruleengine/pkg/idgen"
)

func (s *Service) CreateCondition(input model.Condition) (*model.Condition, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
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
	return c, nil
}

func (s *Service) GetCondition(id string) (*model.Condition, error) {
	return s.store.GetCondition(id)
}

func (s *Service) ListConditions(filter model.ConditionFilter, page, size int) ([]*model.Condition, int, error) {
	all := s.store.ListConditions()
	matched := make([]*model.Condition, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Condition{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateCondition(id string, input model.Condition) (*model.Condition, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	c, err := s.store.GetCondition(id)
	if err != nil {
		return nil, err
	}
	c.Name = input.Name
	c.Field = input.Field
	c.Operator = input.Operator
	c.Value = input.Value
	c.ValueType = input.ValueType
	c.Conjunction = input.Conjunction
	c.UpdatedAt = time.Now()
	if err := s.store.UpdateCondition(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCondition(id string) error {
	return s.store.DeleteCondition(id)
}
