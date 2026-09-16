package service

import (
	"sort"
	"time"

	"ruleengine/internal/model"
	"ruleengine/pkg/idgen"
)

func (s *Service) CreateExecutionLog(input model.ExecutionLog) (*model.ExecutionLog, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	el := &model.ExecutionLog{
		ID:         idgen.Hex(),
		RuleSetID:  input.RuleSetID,
		RuleID:     input.RuleID,
		InputKey:   input.InputKey,
		DurationMs: input.DurationMs,
		Hit:        input.Hit,
		ExecutedAt: now,
	}
	if err := s.store.CreateExecutionLog(el); err != nil {
		return nil, err
	}
	return el, nil
}

func (s *Service) GetExecutionLog(id string) (*model.ExecutionLog, error) {
	return s.store.GetExecutionLog(id)
}

func (s *Service) ListExecutionLogs(filter model.ExecutionLogFilter, page, size int) ([]*model.ExecutionLog, int, error) {
	all := s.store.ListExecutionLogs()
	matched := make([]*model.ExecutionLog, 0, len(all))
	for _, el := range all {
		if filter.Match(el) {
			matched = append(matched, el)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].ExecutedAt.After(matched[j].ExecutedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ExecutionLog{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteExecutionLog(id string) error {
	return s.store.DeleteExecutionLog(id)
}
