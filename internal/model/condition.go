package model

import (
	"strings"
	"time"
)

const (
	OpEQ        = "eq"
	OpNE        = "ne"
	OpGT        = "gt"
	OpGE        = "ge"
	OpLT        = "lt"
	OpLE        = "le"
	OpContains  = "contains"
	OpIn        = "in"
	OpBetween   = "between"
)

const (
	ValueTypeString = "string"
	ValueTypeNumber = "number"
	ValueTypeBool   = "bool"
)

const (
	ConjunctionAnd = "and"
	ConjunctionOr  = "or"
)

var ValidOperators = map[string]bool{
	OpEQ: true, OpNE: true, OpGT: true, OpGE: true,
	OpLT: true, OpLE: true, OpContains: true, OpIn: true, OpBetween: true,
}

type Condition struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Field       string    `json:"field"`
	Operator    string    `json:"operator"`
	Value       string    `json:"value"`
	ValueType   string    `json:"value_type"`
	Conjunction string    `json:"conjunction"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Condition) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	c.Field = strings.TrimSpace(c.Field)
	c.Operator = strings.TrimSpace(c.Operator)
	c.ValueType = strings.TrimSpace(c.ValueType)
	if c.Name == "" {
		return NewValidationError("name", "条件名称不能为空")
	}
	if c.Field == "" {
		return NewValidationError("field", "字段名不能为空")
	}
	if !ValidOperators[c.Operator] {
		return NewValidationError("operator", "操作符不合法")
	}
	if c.ValueType == "" {
		c.ValueType = ValueTypeString
	}
	if c.ValueType != ValueTypeString && c.ValueType != ValueTypeNumber && c.ValueType != ValueTypeBool {
		return NewValidationError("value_type", "值类型不合法")
	}
	if c.Conjunction == "" {
		c.Conjunction = ConjunctionAnd
	}
	if c.Conjunction != ConjunctionAnd && c.Conjunction != ConjunctionOr {
		return NewValidationError("conjunction", "连接词不合法")
	}
	return nil
}

type ConditionFilter struct {
	Name      string
	Field     string
	Operator  string
	ValueType string
	Keyword   string
}

func (f ConditionFilter) Match(c *Condition) bool {
	if f.Operator != "" && c.Operator != f.Operator {
		return false
	}
	if f.ValueType != "" && c.ValueType != f.ValueType {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Name), k) &&
			!strings.Contains(strings.ToLower(c.Field), k) {
			return false
		}
	}
	return true
}
