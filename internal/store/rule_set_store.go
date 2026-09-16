package store

import (
	"strings"

	"ruleengine/internal/model"
)

func (s *MemoryStore) CreateRuleSet(rs *model.RuleSet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.ruleSets {
		if strings.EqualFold(exist.Name, rs.Name) {
			return ErrConflict
		}
	}
	s.ruleSets[rs.ID] = rs
	return nil
}

func (s *MemoryStore) GetRuleSet(id string) (*model.RuleSet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rs, ok := s.ruleSets[id]
	if !ok {
		return nil, ErrNotFound
	}
	return rs, nil
}

func (s *MemoryStore) GetRuleSetByName(name string) (*model.RuleSet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, rs := range s.ruleSets {
		if strings.EqualFold(rs.Name, name) {
			return rs, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListRuleSets() []*model.RuleSet {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RuleSet, 0, len(s.ruleSets))
	for _, rs := range s.ruleSets {
		list = append(list, rs)
	}
	return list
}

func (s *MemoryStore) UpdateRuleSet(rs *model.RuleSet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.ruleSets[rs.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.ruleSets {
		if exist.ID != rs.ID && strings.EqualFold(exist.Name, rs.Name) {
			return ErrConflict
		}
	}
	s.ruleSets[rs.ID] = rs
	return nil
}

func (s *MemoryStore) DeleteRuleSet(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.ruleSets[id]; !ok {
		return ErrNotFound
	}
	delete(s.ruleSets, id)
	return nil
}
