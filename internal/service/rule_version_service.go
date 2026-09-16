package service

import (
	"sort"
	"time"

	"ruleengine/internal/model"
	"ruleengine/pkg/idgen"
)

func (s *Service) CreateRuleVersion(input model.RuleVersion) (*model.RuleVersion, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetRuleSet(input.RuleSetID); err != nil {
		return nil, model.NewValidationError("rule_set_id", "规则集不存在")
	}
	now := time.Now()
	rv := &model.RuleVersion{
		ID:          idgen.Hex(),
		RuleSetID:   input.RuleSetID,
		Version:     input.Version,
		ChangeNote:  input.ChangeNote,
		ChangedBy:   input.ChangedBy,
		PublishedAt: now,
		Status:      input.Status,
		Snapshot:    input.Snapshot,
		CreatedAt:   now,
	}
	if rv.Status == "" {
		rv.Status = model.RuleVersionStatusActive
	}
	if err := s.store.CreateRuleVersion(rv); err != nil {
		return nil, err
	}
	return rv, nil
}

func (s *Service) GetRuleVersion(id string) (*model.RuleVersion, error) {
	return s.store.GetRuleVersion(id)
}

func (s *Service) ListRuleVersions(filter model.RuleVersionFilter, page, size int) ([]*model.RuleVersion, int, error) {
	all := s.store.ListRuleVersions()
	matched := make([]*model.RuleVersion, 0, len(all))
	for _, rv := range all {
		if filter.Match(rv) {
			matched = append(matched, rv)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].Version != matched[j].Version {
			return matched[i].Version > matched[j].Version
		}
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RuleVersion{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRuleVersion(id string, input model.RuleVersion) (*model.RuleVersion, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	rv, err := s.store.GetRuleVersion(id)
	if err != nil {
		return nil, err
	}
	if input.Status != "" && input.Status != rv.Status {
		if !model.RuleVersionCanTransition(rv.Status, input.Status) {
			return nil, model.NewValidationError("status", "版本状态流转不合法")
		}
		rv.Status = input.Status
	}
	rv.ChangeNote = input.ChangeNote
	rv.ChangedBy = input.ChangedBy
	if err := s.store.UpdateRuleVersion(rv); err != nil {
		return nil, err
	}
	return rv, nil
}

func (s *Service) DeleteRuleVersion(id string) error {
	return s.store.DeleteRuleVersion(id)
}

func (s *Service) RollbackRuleVersion(id string, operator string) (*model.RuleSet, error) {
	rv, err := s.store.GetRuleVersion(id)
	if err != nil {
		return nil, err
	}
	if rv.Status != model.RuleVersionStatusActive {
		return nil, model.NewValidationError("status", "只能回滚 active 状态的版本")
	}
	rs, err := s.store.GetRuleSet(rv.RuleSetID)
	if err != nil {
		return nil, err
	}
	rs.Version = rv.Version
	rs.UpdatedAt = time.Now()
	if err := s.store.UpdateRuleSet(rs); err != nil {
		return nil, err
	}
	rv.Status = model.RuleVersionStatusRolledBack
	_ = s.store.UpdateRuleVersion(rv)

	ar := &model.AuditRecord{
		ID:             idgen.Hex(),
		Operator:       operator,
		Operation:      model.AuditOpRollback,
		TargetType:     "rule_set",
		TargetID:       rs.ID,
		BeforeSnapshot: rv.Snapshot,
		AfterSnapshot:  "rolled back",
		OperatedAt:     time.Now(),
	}
	_ = s.store.CreateAuditRecord(ar)
	return rs, nil
}
