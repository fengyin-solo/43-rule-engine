package store

import (
	"testing"

	"ruleengine/internal/model"
	"ruleengine/pkg/idgen"
)

func TestMemoryStore_RuleSet(t *testing.T) {
	s := NewMemoryStore()
	rs := &model.RuleSet{ID: idgen.Hex(), Name: "test", Status: model.RuleSetStatusDraft}
	if err := s.CreateRuleSet(rs); err != nil {
		t.Fatalf("create ruleset: %v", err)
	}
	got, err := s.GetRuleSet(rs.ID)
	if err != nil {
		t.Fatalf("get ruleset: %v", err)
	}
	if got.Name != "test" {
		t.Fatalf("name mismatch")
	}
	if _, err := s.GetRuleSet("not-exist"); err != ErrNotFound {
		t.Fatalf("expected not found")
	}
	rs2 := &model.RuleSet{ID: idgen.Hex(), Name: "test", Status: model.RuleSetStatusDraft}
	if err := s.CreateRuleSet(rs2); err != ErrConflict {
		t.Fatalf("expected conflict")
	}
	if err := s.UpdateRuleSet(rs); err != nil {
		t.Fatalf("update ruleset: %v", err)
	}
	if err := s.DeleteRuleSet(rs.ID); err != nil {
		t.Fatalf("delete ruleset: %v", err)
	}
	if _, err := s.GetRuleSet(rs.ID); err != ErrNotFound {
		t.Fatalf("expected not found after delete")
	}
}

func TestMemoryStore_Rule(t *testing.T) {
	s := NewMemoryStore()
	r := &model.Rule{ID: idgen.Hex(), RuleSetID: idgen.Hex(), Name: "r1", Status: model.RuleStatusActive}
	if err := s.CreateRule(r); err != nil {
		t.Fatalf("create rule: %v", err)
	}
	got, err := s.GetRule(r.ID)
	if err != nil {
		t.Fatalf("get rule: %v", err)
	}
	if got.Name != "r1" {
		t.Fatalf("name mismatch")
	}
	if err := s.DeleteRule(r.ID); err != nil {
		t.Fatalf("delete rule: %v", err)
	}
}

func TestMemoryStore_Condition(t *testing.T) {
	s := NewMemoryStore()
	c := &model.Condition{ID: idgen.Hex(), Name: "c1", Field: "age", Operator: model.OpEQ, Value: "18", ValueType: model.ValueTypeNumber}
	if err := s.CreateCondition(c); err != nil {
		t.Fatalf("create condition: %v", err)
	}
	got, err := s.GetCondition(c.ID)
	if err != nil {
		t.Fatalf("get condition: %v", err)
	}
	if got.Field != "age" {
		t.Fatalf("field mismatch")
	}
	if err := s.DeleteCondition(c.ID); err != nil {
		t.Fatalf("delete condition: %v", err)
	}
}

func TestMemoryStore_Action(t *testing.T) {
	s := NewMemoryStore()
	a := &model.Action{ID: idgen.Hex(), Name: "a1", ActionType: model.ActionTypeLog}
	if err := s.CreateAction(a); err != nil {
		t.Fatalf("create action: %v", err)
	}
	got, err := s.GetAction(a.ID)
	if err != nil {
		t.Fatalf("get action: %v", err)
	}
	if got.Name != "a1" {
		t.Fatalf("name mismatch")
	}
	if err := s.DeleteAction(a.ID); err != nil {
		t.Fatalf("delete action: %v", err)
	}
}

func TestMemoryStore_EvaluationRecord(t *testing.T) {
	s := NewMemoryStore()
	e := &model.EvaluationRecord{ID: idgen.Hex(), RuleSetID: idgen.Hex(), Result: model.ResultHit}
	if err := s.CreateEvaluationRecord(e); err != nil {
		t.Fatalf("create eval: %v", err)
	}
	got, err := s.GetEvaluationRecord(e.ID)
	if err != nil {
		t.Fatalf("get eval: %v", err)
	}
	if got.Result != model.ResultHit {
		t.Fatalf("result mismatch")
	}
	if err := s.DeleteEvaluationRecord(e.ID); err != nil {
		t.Fatalf("delete eval: %v", err)
	}
}

func TestMemoryStore_RuleVersion(t *testing.T) {
	s := NewMemoryStore()
	rv := &model.RuleVersion{ID: idgen.Hex(), RuleSetID: idgen.Hex(), Version: 1, Status: model.RuleVersionStatusActive}
	if err := s.CreateRuleVersion(rv); err != nil {
		t.Fatalf("create version: %v", err)
	}
	got, err := s.GetRuleVersion(rv.ID)
	if err != nil {
		t.Fatalf("get version: %v", err)
	}
	if got.Version != 1 {
		t.Fatalf("version mismatch")
	}
	if err := s.DeleteRuleVersion(rv.ID); err != nil {
		t.Fatalf("delete version: %v", err)
	}
}

func TestMemoryStore_ExecutionLog(t *testing.T) {
	s := NewMemoryStore()
	el := &model.ExecutionLog{ID: idgen.Hex(), RuleSetID: idgen.Hex(), RuleID: idgen.Hex(), DurationMs: 10, Hit: true}
	if err := s.CreateExecutionLog(el); err != nil {
		t.Fatalf("create log: %v", err)
	}
	got, err := s.GetExecutionLog(el.ID)
	if err != nil {
		t.Fatalf("get log: %v", err)
	}
	if !got.Hit {
		t.Fatalf("hit mismatch")
	}
	if err := s.DeleteExecutionLog(el.ID); err != nil {
		t.Fatalf("delete log: %v", err)
	}
}

func TestMemoryStore_DataType(t *testing.T) {
	s := NewMemoryStore()
	dt := &model.DataType{ID: idgen.Hex(), Name: "user", Fields: []byte("{}"), Enabled: true}
	if err := s.CreateDataType(dt); err != nil {
		t.Fatalf("create datatype: %v", err)
	}
	got, err := s.GetDataType(dt.ID)
	if err != nil {
		t.Fatalf("get datatype: %v", err)
	}
	if got.Name != "user" {
		t.Fatalf("name mismatch")
	}
	if err := s.DeleteDataType(dt.ID); err != nil {
		t.Fatalf("delete datatype: %v", err)
	}
}

func TestMemoryStore_AuditRecord(t *testing.T) {
	s := NewMemoryStore()
	ar := &model.AuditRecord{ID: idgen.Hex(), Operator: "admin", Operation: model.AuditOpCreate, TargetType: "rule_set", TargetID: idgen.Hex()}
	if err := s.CreateAuditRecord(ar); err != nil {
		t.Fatalf("create audit: %v", err)
	}
	got, err := s.GetAuditRecord(ar.ID)
	if err != nil {
		t.Fatalf("get audit: %v", err)
	}
	if got.Operator != "admin" {
		t.Fatalf("operator mismatch")
	}
	if err := s.DeleteAuditRecord(ar.ID); err != nil {
		t.Fatalf("delete audit: %v", err)
	}
}

func TestMemoryStore_RuleStat(t *testing.T) {
	s := NewMemoryStore()
	rs := &model.RuleStat{ID: idgen.Hex(), RuleSetID: idgen.Hex(), RuleID: idgen.Hex(), Date: "2024-01-01", TotalEval: 10, HitCount: 5}
	if err := s.CreateRuleStat(rs); err != nil {
		t.Fatalf("create stat: %v", err)
	}
	got, err := s.GetRuleStat(rs.ID)
	if err != nil {
		t.Fatalf("get stat: %v", err)
	}
	if got.HitCount != 5 {
		t.Fatalf("hit count mismatch")
	}
	rs2 := &model.RuleStat{ID: idgen.Hex(), RuleSetID: rs.RuleSetID, RuleID: rs.RuleID, Date: rs.Date, TotalEval: 1}
	if err := s.CreateRuleStat(rs2); err != ErrConflict {
		t.Fatalf("expected conflict")
	}
	if err := s.DeleteRuleStat(rs.ID); err != nil {
		t.Fatalf("delete stat: %v", err)
	}
}
