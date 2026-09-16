package service

import (
	"sort"
	"time"

	"ruleengine/internal/model"
	"ruleengine/pkg/idgen"
)

func (s *Service) CreateEvaluationRecord(input model.EvaluationRecord) (*model.EvaluationRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	er := &model.EvaluationRecord{
		ID:             idgen.Hex(),
		RuleSetID:      input.RuleSetID,
		InputSnapshot:  input.InputSnapshot,
		MatchedRuleIDs: input.MatchedRuleIDs,
		Result:         input.Result,
		EvaluatedAt:    now,
	}
	if err := s.store.CreateEvaluationRecord(er); err != nil {
		return nil, err
	}
	return er, nil
}

func (s *Service) GetEvaluationRecord(id string) (*model.EvaluationRecord, error) {
	return s.store.GetEvaluationRecord(id)
}

func (s *Service) ListEvaluationRecords(filter model.EvaluationRecordFilter, page, size int) ([]*model.EvaluationRecord, int, error) {
	all := s.store.ListEvaluationRecords()
	matched := make([]*model.EvaluationRecord, 0, len(all))
	for _, er := range all {
		if filter.Match(er) {
			matched = append(matched, er)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].EvaluatedAt.After(matched[j].EvaluatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.EvaluationRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteEvaluationRecord(id string) error {
	return s.store.DeleteEvaluationRecord(id)
}
