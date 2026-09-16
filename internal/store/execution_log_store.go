package store

import (
	"ruleengine/internal/model"
)

func (s *MemoryStore) CreateExecutionLog(el *model.ExecutionLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.executionLogs[el.ID] = el
	return nil
}

func (s *MemoryStore) GetExecutionLog(id string) (*model.ExecutionLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	el, ok := s.executionLogs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return el, nil
}

func (s *MemoryStore) ListExecutionLogs() []*model.ExecutionLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ExecutionLog, 0, len(s.executionLogs))
	for _, el := range s.executionLogs {
		list = append(list, el)
	}
	return list
}

func (s *MemoryStore) UpdateExecutionLog(el *model.ExecutionLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.executionLogs[el.ID]; !ok {
		return ErrNotFound
	}
	s.executionLogs[el.ID] = el
	return nil
}

func (s *MemoryStore) DeleteExecutionLog(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.executionLogs[id]; !ok {
		return ErrNotFound
	}
	delete(s.executionLogs, id)
	return nil
}
