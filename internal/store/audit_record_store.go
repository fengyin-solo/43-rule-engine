package store

import (
	"ruleengine/internal/model"
)

func (s *MemoryStore) CreateAuditRecord(ar *model.AuditRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.auditRecords[ar.ID] = ar
	return nil
}

func (s *MemoryStore) GetAuditRecord(id string) (*model.AuditRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ar, ok := s.auditRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return ar, nil
}

func (s *MemoryStore) ListAuditRecords() []*model.AuditRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.AuditRecord, 0, len(s.auditRecords))
	for _, ar := range s.auditRecords {
		list = append(list, ar)
	}
	return list
}

func (s *MemoryStore) UpdateAuditRecord(ar *model.AuditRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.auditRecords[ar.ID]; !ok {
		return ErrNotFound
	}
	s.auditRecords[ar.ID] = ar
	return nil
}

func (s *MemoryStore) DeleteAuditRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.auditRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.auditRecords, id)
	return nil
}
