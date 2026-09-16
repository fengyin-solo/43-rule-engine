package service

import (
	"sort"
	"time"

	"ruleengine/internal/model"
	"ruleengine/pkg/idgen"
)

func (s *Service) CreateAuditRecord(input model.AuditRecord) (*model.AuditRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	ar := &model.AuditRecord{
		ID:             idgen.Hex(),
		Operator:       input.Operator,
		Operation:      input.Operation,
		TargetType:     input.TargetType,
		TargetID:       input.TargetID,
		BeforeSnapshot: input.BeforeSnapshot,
		AfterSnapshot:  input.AfterSnapshot,
		OperatedAt:     now,
	}
	if err := s.store.CreateAuditRecord(ar); err != nil {
		return nil, err
	}
	return ar, nil
}

func (s *Service) GetAuditRecord(id string) (*model.AuditRecord, error) {
	return s.store.GetAuditRecord(id)
}

func (s *Service) ListAuditRecords(filter model.AuditRecordFilter, page, size int) ([]*model.AuditRecord, int, error) {
	all := s.store.ListAuditRecords()
	matched := make([]*model.AuditRecord, 0, len(all))
	for _, ar := range all {
		if filter.Match(ar) {
			matched = append(matched, ar)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].OperatedAt.After(matched[j].OperatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.AuditRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteAuditRecord(id string) error {
	return s.store.DeleteAuditRecord(id)
}
