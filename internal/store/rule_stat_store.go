package store

import (
	"ruleengine/internal/model"
)

func (s *MemoryStore) CreateRuleStat(rs *model.RuleStat) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.ruleStats {
		if exist.RuleSetID == rs.RuleSetID && exist.RuleID == rs.RuleID && exist.Date == rs.Date {
			return ErrConflict
		}
	}
	s.ruleStats[rs.ID] = rs
	return nil
}

func (s *MemoryStore) GetRuleStat(id string) (*model.RuleStat, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rs, ok := s.ruleStats[id]
	if !ok {
		return nil, ErrNotFound
	}
	return rs, nil
}

func (s *MemoryStore) ListRuleStats() []*model.RuleStat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RuleStat, 0, len(s.ruleStats))
	for _, rs := range s.ruleStats {
		list = append(list, rs)
	}
	return list
}

func (s *MemoryStore) UpdateRuleStat(rs *model.RuleStat) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.ruleStats[rs.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.ruleStats {
		if exist.ID != rs.ID && exist.RuleSetID == rs.RuleSetID && exist.RuleID == rs.RuleID && exist.Date == rs.Date {
			return ErrConflict
		}
	}
	s.ruleStats[rs.ID] = rs
	return nil
}

func (s *MemoryStore) DeleteRuleStat(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.ruleStats[id]; !ok {
		return ErrNotFound
	}
	delete(s.ruleStats, id)
	return nil
}
