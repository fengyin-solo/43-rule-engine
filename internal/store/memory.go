package store

import (
	"sync"

	"ruleengine/internal/model"
)

type MemoryStore struct {
	mu                 sync.RWMutex
	ruleSets           map[string]*model.RuleSet
	rules              map[string]*model.Rule
	conditions         map[string]*model.Condition
	actions            map[string]*model.Action
	evaluationRecords  map[string]*model.EvaluationRecord
	ruleVersions       map[string]*model.RuleVersion
	executionLogs      map[string]*model.ExecutionLog
	dataTypes          map[string]*model.DataType
	auditRecords       map[string]*model.AuditRecord
	ruleStats          map[string]*model.RuleStat
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		ruleSets:          make(map[string]*model.RuleSet),
		rules:             make(map[string]*model.Rule),
		conditions:        make(map[string]*model.Condition),
		actions:           make(map[string]*model.Action),
		evaluationRecords: make(map[string]*model.EvaluationRecord),
		ruleVersions:      make(map[string]*model.RuleVersion),
		executionLogs:     make(map[string]*model.ExecutionLog),
		dataTypes:         make(map[string]*model.DataType),
		auditRecords:      make(map[string]*model.AuditRecord),
		ruleStats:         make(map[string]*model.RuleStat),
	}
}

var _ Store = (*MemoryStore)(nil)
