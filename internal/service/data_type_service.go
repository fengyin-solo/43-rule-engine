package service

import (
	"sort"
	"time"

	"ruleengine/internal/model"
	"ruleengine/pkg/idgen"
)

func (s *Service) CreateDataType(input model.DataType) (*model.DataType, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	dt := &model.DataType{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Fields:      input.Fields,
		Description: input.Description,
		Enabled:     input.Enabled,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateDataType(dt); err != nil {
		return nil, err
	}
	return dt, nil
}

func (s *Service) GetDataType(id string) (*model.DataType, error) {
	return s.store.GetDataType(id)
}

func (s *Service) ListDataTypes(filter model.DataTypeFilter, page, size int) ([]*model.DataType, int, error) {
	all := s.store.ListDataTypes()
	matched := make([]*model.DataType, 0, len(all))
	for _, dt := range all {
		if filter.Match(dt) {
			matched = append(matched, dt)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.DataType{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateDataType(id string, input model.DataType) (*model.DataType, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	dt, err := s.store.GetDataType(id)
	if err != nil {
		return nil, err
	}
	dt.Name = input.Name
	dt.Fields = input.Fields
	dt.Description = input.Description
	dt.Enabled = input.Enabled
	dt.UpdatedAt = time.Now()
	if err := s.store.UpdateDataType(dt); err != nil {
		return nil, err
	}
	return dt, nil
}

func (s *Service) DeleteDataType(id string) error {
	return s.store.DeleteDataType(id)
}
