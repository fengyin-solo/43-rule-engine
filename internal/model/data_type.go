package model

import (
	"encoding/json"
	"strings"
	"time"
)

type DataType struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Fields      json.RawMessage `json:"fields"`
	Description string          `json:"description"`
	Enabled     bool            `json:"enabled"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (dt *DataType) Validate() error {
	dt.Name = strings.TrimSpace(dt.Name)
	dt.Description = strings.TrimSpace(dt.Description)
	if dt.Name == "" {
		return NewValidationError("name", "数据类型名称不能为空")
	}
	if len(dt.Fields) == 0 {
		dt.Fields = json.RawMessage("{}")
	}
	return nil
}

type DataTypeFilter struct {
	Name    string
	Enabled *bool
	Keyword string
}

func (f DataTypeFilter) Match(dt *DataType) bool {
	if f.Enabled != nil && dt.Enabled != *f.Enabled {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(dt.Name), k) &&
			!strings.Contains(strings.ToLower(dt.Description), k) {
			return false
		}
	}
	return true
}
