package store

import (
	"ruleengine/internal/model"
)

func (s *MemoryStore) CreateEvaluationRecord(er *model.EvaluationRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evaluationRecords[er.ID] = er
	return nil
}

func (s *MemoryStore) GetEvaluationRecord(id string) (*model.EvaluationRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	er, ok := s.evaluationRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return er, nil
}

func (s *MemoryStore) ListEvaluationRecords() []*model.EvaluationRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.EvaluationRecord, 0, len(s.evaluationRecords))
	for _, er := range s.evaluationRecords {
		list = append(list, er)
	}
	return list
}

func (s *MemoryStore) UpdateEvaluationRecord(er *model.EvaluationRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.evaluationRecords[er.ID]; !ok {
		return ErrNotFound
	}
	s.evaluationRecords[er.ID] = er
	return nil
}

func (s *MemoryStore) DeleteEvaluationRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.evaluationRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.evaluationRecords, id)
	return nil
}
