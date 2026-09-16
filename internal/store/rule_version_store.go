package store

import (
	"ruleengine/internal/model"
)

func (s *MemoryStore) CreateRuleVersion(rv *model.RuleVersion) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ruleVersions[rv.ID] = rv
	return nil
}

func (s *MemoryStore) GetRuleVersion(id string) (*model.RuleVersion, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rv, ok := s.ruleVersions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return rv, nil
}

func (s *MemoryStore) ListRuleVersions() []*model.RuleVersion {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RuleVersion, 0, len(s.ruleVersions))
	for _, rv := range s.ruleVersions {
		list = append(list, rv)
	}
	return list
}

func (s *MemoryStore) ListRuleVersionsByRuleSetID(ruleSetID string) []*model.RuleVersion {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RuleVersion, 0)
	for _, rv := range s.ruleVersions {
		if rv.RuleSetID == ruleSetID {
			list = append(list, rv)
		}
	}
	return list
}

func (s *MemoryStore) UpdateRuleVersion(rv *model.RuleVersion) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.ruleVersions[rv.ID]; !ok {
		return ErrNotFound
	}
	s.ruleVersions[rv.ID] = rv
	return nil
}

func (s *MemoryStore) DeleteRuleVersion(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.ruleVersions[id]; !ok {
		return ErrNotFound
	}
	delete(s.ruleVersions, id)
	return nil
}
