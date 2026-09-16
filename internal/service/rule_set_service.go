package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"ruleengine/internal/model"
	"ruleengine/pkg/idgen"
)

func (s *Service) CreateRuleSet(input model.RuleSet) (*model.RuleSet, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	rs := &model.RuleSet{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Description: input.Description,
		Version:     input.Version,
		Status:      input.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateRuleSet(rs); err != nil {
		return nil, err
	}
	s.log.Infof("创建规则集 %s", rs.ID)
	return rs, nil
}

func (s *Service) GetRuleSet(id string) (*model.RuleSet, error) {
	return s.store.GetRuleSet(id)
}

func (s *Service) ListRuleSets(filter model.RuleSetFilter, page, size int) ([]*model.RuleSet, int, error) {
	all := s.store.ListRuleSets()
	matched := make([]*model.RuleSet, 0, len(all))
	for _, rs := range all {
		if filter.Match(rs) {
			matched = append(matched, rs)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RuleSet{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRuleSet(id string, input model.RuleSet) (*model.RuleSet, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	rs, err := s.store.GetRuleSet(id)
	if err != nil {
		return nil, err
	}
	if rs.Status == model.RuleSetStatusDisabled {
		return nil, model.NewValidationError("status", "已禁用的规则集不能修改")
	}
	if rs.Status == model.RuleSetStatusPublished {
		return nil, model.NewValidationError("status", "已发布的规则集不能直接修改，请先禁用")
	}
	rs.Name = input.Name
	rs.Description = input.Description
	rs.UpdatedAt = time.Now()
	if err := s.store.UpdateRuleSet(rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (s *Service) DeleteRuleSet(id string) error {
	rs, err := s.store.GetRuleSet(id)
	if err != nil {
		return err
	}
	if rs.Status == model.RuleSetStatusPublished {
		return model.NewValidationError("status", "已发布的规则集不能删除")
	}
	return s.store.DeleteRuleSet(id)
}

func (s *Service) PublishRuleSet(id string, changedBy string) (*model.RuleSet, error) {
	rs, err := s.store.GetRuleSet(id)
	if err != nil {
		return nil, err
	}
	if !model.RuleSetCanTransition(rs.Status, model.RuleSetStatusPublished) {
		return nil, model.NewValidationError("status", fmt.Sprintf("规则集当前状态 %s 不能发布", rs.Status))
	}
	oldStatus := rs.Status
	rs.Status = model.RuleSetStatusPublished
	rs.Version++
	rs.UpdatedAt = time.Now()

	snapshot, _ := json.Marshal(rs)
	rv := &model.RuleVersion{
		ID:          idgen.Hex(),
		RuleSetID:   rs.ID,
		Version:     rs.Version,
		ChangeNote:  fmt.Sprintf("从 %s 发布到 published", oldStatus),
		ChangedBy:   changedBy,
		PublishedAt: time.Now(),
		Status:      model.RuleVersionStatusActive,
		Snapshot:    string(snapshot),
		CreatedAt:   time.Now(),
	}
	if err := rv.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateRuleVersion(rv); err != nil {
		return nil, err
	}

	ar := &model.AuditRecord{
		ID:             idgen.Hex(),
		Operator:       changedBy,
		Operation:      model.AuditOpPublish,
		TargetType:     "rule_set",
		TargetID:       rs.ID,
		BeforeSnapshot: oldStatus,
		AfterSnapshot:  rs.Status,
		OperatedAt:     time.Now(),
	}
	if err := ar.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateAuditRecord(ar); err != nil {
		return nil, err
	}

	if err := s.store.UpdateRuleSet(rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (s *Service) DisableRuleSet(id string, changedBy string) (*model.RuleSet, error) {
	rs, err := s.store.GetRuleSet(id)
	if err != nil {
		return nil, err
	}
	if !model.RuleSetCanTransition(rs.Status, model.RuleSetStatusDisabled) {
		return nil, model.NewValidationError("status", fmt.Sprintf("规则集当前状态 %s 不能禁用", rs.Status))
	}
	oldStatus := rs.Status
	rs.Status = model.RuleSetStatusDisabled
	rs.UpdatedAt = time.Now()

	ar := &model.AuditRecord{
		ID:             idgen.Hex(),
		Operator:       changedBy,
		Operation:      model.AuditOpUpdate,
		TargetType:     "rule_set",
		TargetID:       rs.ID,
		BeforeSnapshot: oldStatus,
		AfterSnapshot:  rs.Status,
		OperatedAt:     time.Now(),
	}
	if err := ar.Validate(); err != nil {
		return nil, err
	}
	_ = s.store.CreateAuditRecord(ar)

	if err := s.store.UpdateRuleSet(rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (s *Service) ExportRuleSetSnapshot(id string) (string, error) {
	rs, err := s.store.GetRuleSet(id)
	if err != nil {
		return "", err
	}
	rules := s.store.ListRulesByRuleSetID(id)
	condMap := make(map[string]*model.Condition)
	actMap := make(map[string]*model.Action)
	for _, r := range rules {
		for _, cid := range r.ConditionIDs {
			if c, err := s.store.GetCondition(cid); err == nil {
				condMap[cid] = c
			}
		}
		for _, aid := range r.ActionIDs {
			if a, err := s.store.GetAction(aid); err == nil {
				actMap[aid] = a
			}
		}
	}
	snapshot := map[string]interface{}{
		"rule_set":   rs,
		"rules":      rules,
		"conditions": condMap,
		"actions":    actMap,
	}
	b, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *Service) ImportRuleSetSnapshot(data string, operator string) (*model.RuleSet, error) {
	var payload struct {
		RuleSet    *model.RuleSet            `json:"rule_set"`
		Rules      []*model.Rule             `json:"rules"`
		Conditions map[string]*model.Condition `json:"conditions"`
		Actions    map[string]*model.Action  `json:"actions"`
	}
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return nil, model.NewValidationError("data", "快照解析失败: "+err.Error())
	}
	if payload.RuleSet == nil {
		return nil, model.NewValidationError("rule_set", "快照缺少规则集")
	}
	oldID := payload.RuleSet.ID
	payload.RuleSet.ID = idgen.Hex()
	payload.RuleSet.Name = payload.RuleSet.Name + "-imported"
	payload.RuleSet.Status = model.RuleSetStatusDraft
	payload.RuleSet.Version = 1
	payload.RuleSet.CreatedAt = time.Now()
	payload.RuleSet.UpdatedAt = time.Now()
	if err := payload.RuleSet.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateRuleSet(payload.RuleSet); err != nil {
		return nil, err
	}
	idMapping := make(map[string]string)
	idMapping[oldID] = payload.RuleSet.ID
	for _, c := range payload.Conditions {
		oldCID := c.ID
		c.ID = idgen.Hex()
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
		if err := c.Validate(); err == nil {
			_ = s.store.CreateCondition(c)
			idMapping[oldCID] = c.ID
		}
	}
	for _, a := range payload.Actions {
		oldAID := a.ID
		a.ID = idgen.Hex()
		a.CreatedAt = time.Now()
		a.UpdatedAt = time.Now()
		if err := a.Validate(); err == nil {
			_ = s.store.CreateAction(a)
			idMapping[oldAID] = a.ID
		}
	}
	for _, r := range payload.Rules {
		r.ID = idgen.Hex()
		r.RuleSetID = payload.RuleSet.ID
		for i := range r.ConditionIDs {
			if nid, ok := idMapping[r.ConditionIDs[i]]; ok {
				r.ConditionIDs[i] = nid
			}
		}
		for i := range r.ActionIDs {
			if nid, ok := idMapping[r.ActionIDs[i]]; ok {
				r.ActionIDs[i] = nid
			}
		}
		r.CreatedAt = time.Now()
		r.UpdatedAt = time.Now()
		if err := r.Validate(); err == nil {
			_ = s.store.CreateRule(r)
		}
	}
	ar := &model.AuditRecord{
		ID:            idgen.Hex(),
		Operator:      operator,
		Operation:     model.AuditOpCreate,
		TargetType:    "rule_set",
		TargetID:      payload.RuleSet.ID,
		AfterSnapshot: fmt.Sprintf("imported from %s", oldID),
		OperatedAt:    time.Now(),
	}
	_ = s.store.CreateAuditRecord(ar)
	return payload.RuleSet, nil
}
