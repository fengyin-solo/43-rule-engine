package store

import (
	"strings"

	"ruleengine/internal/model"
)

func (s *MemoryStore) CreateAction(a *model.Action) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.actions {
		if strings.EqualFold(exist.Name, a.Name) {
			return ErrConflict
		}
	}
	s.actions[a.ID] = a
	return nil
}

func (s *MemoryStore) GetAction(id string) (*model.Action, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.actions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *MemoryStore) ListActions() []*model.Action {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Action, 0, len(s.actions))
	for _, a := range s.actions {
		list = append(list, a)
	}
	return list
}

func (s *MemoryStore) UpdateAction(a *model.Action) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.actions[a.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.actions {
		if exist.ID != a.ID && strings.EqualFold(exist.Name, a.Name) {
			return ErrConflict
		}
	}
	s.actions[a.ID] = a
	return nil
}

func (s *MemoryStore) DeleteAction(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.actions[id]; !ok {
		return ErrNotFound
	}
	delete(s.actions, id)
	return nil
}
