package model

import (
	"strings"
	"time"
)

const (
	ActionTypeSetField  = "set_field"
	ActionTypeSendNotify = "send_notify"
	ActionTypeLog       = "log"
	ActionTypeIncrement = "increment"
)

type Action struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	ActionType string            `json:"action_type"`
	Params     map[string]string `json:"params"`
	Target     string            `json:"target"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

func (a *Action) Validate() error {
	a.Name = strings.TrimSpace(a.Name)
	a.ActionType = strings.TrimSpace(a.ActionType)
	a.Target = strings.TrimSpace(a.Target)
	if a.Name == "" {
		return NewValidationError("name", "动作名称不能为空")
	}
	if a.ActionType == "" {
		return NewValidationError("action_type", "动作类型不能为空")
	}
	if a.ActionType != ActionTypeSetField && a.ActionType != ActionTypeSendNotify &&
		a.ActionType != ActionTypeLog && a.ActionType != ActionTypeIncrement {
		return NewValidationError("action_type", "动作类型不合法")
	}
	if a.Params == nil {
		a.Params = make(map[string]string)
	}
	return nil
}

type ActionFilter struct {
	Name       string
	ActionType string
	Keyword    string
}

func (f ActionFilter) Match(a *Action) bool {
	if f.ActionType != "" && a.ActionType != f.ActionType {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(a.Name), k) {
			return false
		}
	}
	return true
}
