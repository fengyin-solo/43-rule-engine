package model

import (
	"strings"
	"time"
)

const (
	AuditOpCreate   = "create"
	AuditOpUpdate   = "update"
	AuditOpPublish  = "publish"
	AuditOpRollback = "rollback"
	AuditOpDelete   = "delete"
)

type AuditRecord struct {
	ID            string    `json:"id"`
	Operator      string    `json:"operator"`
	Operation     string    `json:"operation"`
	TargetType    string    `json:"target_type"`
	TargetID      string    `json:"target_id"`
	BeforeSnapshot string   `json:"before_snapshot"`
	AfterSnapshot  string   `json:"after_snapshot"`
	OperatedAt    time.Time `json:"operated_at"`
}

func (ar *AuditRecord) Validate() error {
	ar.Operator = strings.TrimSpace(ar.Operator)
	ar.Operation = strings.TrimSpace(ar.Operation)
	ar.TargetType = strings.TrimSpace(ar.TargetType)
	if ar.Operator == "" {
		return NewValidationError("operator", "操作人不能为空")
	}
	if ar.Operation == "" {
		return NewValidationError("operation", "操作类型不能为空")
	}
	if ar.Operation != AuditOpCreate && ar.Operation != AuditOpUpdate &&
		ar.Operation != AuditOpPublish && ar.Operation != AuditOpRollback && ar.Operation != AuditOpDelete {
		return NewValidationError("operation", "操作类型不合法")
	}
	if ar.TargetType == "" {
		return NewValidationError("target_type", "目标类型不能为空")
	}
	if ar.TargetID == "" {
		return NewValidationError("target_id", "目标 ID 不能为空")
	}
	return nil
}

type AuditRecordFilter struct {
	Operator   string
	Operation  string
	TargetType string
	TargetID   string
	Keyword    string
}

func (f AuditRecordFilter) Match(ar *AuditRecord) bool {
	if f.Operator != "" && ar.Operator != f.Operator {
		return false
	}
	if f.Operation != "" && ar.Operation != f.Operation {
		return false
	}
	if f.TargetType != "" && ar.TargetType != f.TargetType {
		return false
	}
	if f.TargetID != "" && ar.TargetID != f.TargetID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(ar.Operator), k) {
			return false
		}
	}
	return true
}
