package store

import (
	"strings"

	"ruleengine/internal/model"
)

func (s *MemoryStore) CreateDataType(dt *model.DataType) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.dataTypes {
		if strings.EqualFold(exist.Name, dt.Name) {
			return ErrConflict
		}
	}
	s.dataTypes[dt.ID] = dt
	return nil
}

func (s *MemoryStore) GetDataType(id string) (*model.DataType, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dt, ok := s.dataTypes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return dt, nil
}

func (s *MemoryStore) GetDataTypeByName(name string) (*model.DataType, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, dt := range s.dataTypes {
		if strings.EqualFold(dt.Name, name) {
			return dt, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListDataTypes() []*model.DataType {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.DataType, 0, len(s.dataTypes))
	for _, dt := range s.dataTypes {
		list = append(list, dt)
	}
	return list
}

func (s *MemoryStore) UpdateDataType(dt *model.DataType) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dataTypes[dt.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.dataTypes {
		if exist.ID != dt.ID && strings.EqualFold(exist.Name, dt.Name) {
			return ErrConflict
		}
	}
	s.dataTypes[dt.ID] = dt
	return nil
}

func (s *MemoryStore) DeleteDataType(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dataTypes[id]; !ok {
		return ErrNotFound
	}
	delete(s.dataTypes, id)
	return nil
}
