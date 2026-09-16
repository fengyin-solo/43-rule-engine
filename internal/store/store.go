// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"ruleengine/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// RuleSet
	CreateRuleSet(rs *model.RuleSet) error
	GetRuleSet(id string) (*model.RuleSet, error)
	GetRuleSetByName(name string) (*model.RuleSet, error)
	ListRuleSets() []*model.RuleSet
	UpdateRuleSet(rs *model.RuleSet) error
	DeleteRuleSet(id string) error

	// Rule
	CreateRule(r *model.Rule) error
	GetRule(id string) (*model.Rule, error)
	ListRules() []*model.Rule
	ListRulesByRuleSetID(ruleSetID string) []*model.Rule
	UpdateRule(r *model.Rule) error
	DeleteRule(id string) error

	// Condition
	CreateCondition(c *model.Condition) error
	GetCondition(id string) (*model.Condition, error)
	ListConditions() []*model.Condition
	UpdateCondition(c *model.Condition) error
	DeleteCondition(id string) error

	// Action
	CreateAction(a *model.Action) error
	GetAction(id string) (*model.Action, error)
	ListActions() []*model.Action
	UpdateAction(a *model.Action) error
	DeleteAction(id string) error

	// EvaluationRecord
	CreateEvaluationRecord(er *model.EvaluationRecord) error
	GetEvaluationRecord(id string) (*model.EvaluationRecord, error)
	ListEvaluationRecords() []*model.EvaluationRecord
	UpdateEvaluationRecord(er *model.EvaluationRecord) error
	DeleteEvaluationRecord(id string) error

	// RuleVersion
	CreateRuleVersion(rv *model.RuleVersion) error
	GetRuleVersion(id string) (*model.RuleVersion, error)
	ListRuleVersions() []*model.RuleVersion
	ListRuleVersionsByRuleSetID(ruleSetID string) []*model.RuleVersion
	UpdateRuleVersion(rv *model.RuleVersion) error
	DeleteRuleVersion(id string) error

	// ExecutionLog
	CreateExecutionLog(el *model.ExecutionLog) error
	GetExecutionLog(id string) (*model.ExecutionLog, error)
	ListExecutionLogs() []*model.ExecutionLog
	UpdateExecutionLog(el *model.ExecutionLog) error
	DeleteExecutionLog(id string) error

	// DataType
	CreateDataType(dt *model.DataType) error
	GetDataType(id string) (*model.DataType, error)
	GetDataTypeByName(name string) (*model.DataType, error)
	ListDataTypes() []*model.DataType
	UpdateDataType(dt *model.DataType) error
	DeleteDataType(id string) error

	// AuditRecord
	CreateAuditRecord(ar *model.AuditRecord) error
	GetAuditRecord(id string) (*model.AuditRecord, error)
	ListAuditRecords() []*model.AuditRecord
	UpdateAuditRecord(ar *model.AuditRecord) error
	DeleteAuditRecord(id string) error

	// RuleStat
	CreateRuleStat(rs *model.RuleStat) error
	GetRuleStat(id string) (*model.RuleStat, error)
	ListRuleStats() []*model.RuleStat
	UpdateRuleStat(rs *model.RuleStat) error
	DeleteRuleStat(id string) error
}
