package service

import (
	"sort"
	"time"

	"ruleengine/internal/model"
	"ruleengine/pkg/idgen"
)

func (s *Service) CreateAction(input model.Action) (*model.Action, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
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
	return a, nil
}

func (s *Service) GetAction(id string) (*model.Action, error) {
	return s.store.GetAction(id)
}

func (s *Service) ListActions(filter model.ActionFilter, page, size int) ([]*model.Action, int, error) {
	all := s.store.ListActions()
	matched := make([]*model.Action, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Action{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateAction(id string, input model.Action) (*model.Action, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	a, err := s.store.GetAction(id)
	if err != nil {
		return nil, err
	}
	a.Name = input.Name
	a.ActionType = input.ActionType
	a.Params = input.Params
	a.Target = input.Target
	a.UpdatedAt = time.Now()
	if err := s.store.UpdateAction(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) DeleteAction(id string) error {
	return s.store.DeleteAction(id)
}
