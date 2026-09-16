package service

import (
	"fmt"
	"testing"

	"ruleengine/internal/config"
	"ruleengine/internal/model"
	"ruleengine/internal/store"
	"ruleengine/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100, APIKey: "test-key"}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func TestService_RuleSetLifecycle(t *testing.T) {
	svc := newTestService()
	rs, err := svc.CreateRuleSet(model.RuleSet{Name: "lifecycle", Description: "desc"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if rs.Status != model.RuleSetStatusDraft {
		t.Fatalf("expected draft")
	}
	_, err = svc.PublishRuleSet(rs.ID, "admin")
	if err != nil {
		t.Fatalf("publish: %v", err)
	}
	got, _ := svc.GetRuleSet(rs.ID)
	if got.Status != model.RuleSetStatusPublished {
		t.Fatalf("expected published")
	}
	_, err = svc.UpdateRuleSet(rs.ID, model.RuleSet{Name: "updated"})
	if err == nil {
		t.Fatalf("expected error updating published ruleset")
	}
	_, err = svc.DisableRuleSet(rs.ID, "admin")
	if err != nil {
		t.Fatalf("disable: %v", err)
	}
	got, _ = svc.GetRuleSet(rs.ID)
	if got.Status != model.RuleSetStatusDisabled {
		t.Fatalf("expected disabled")
	}
}

func TestService_RuleSetStatusMachine(t *testing.T) {
	svc := newTestService()
	rs, _ := svc.CreateRuleSet(model.RuleSet{Name: "sm"})
	_, err := svc.DisableRuleSet(rs.ID, "admin")
	if err == nil {
		t.Fatalf("draft cannot disable")
	}
	_, _ = svc.PublishRuleSet(rs.ID, "admin")
	_, err = svc.PublishRuleSet(rs.ID, "admin")
	if err == nil {
		t.Fatalf("published cannot publish again")
	}
}

func TestService_RuleCRUD(t *testing.T) {
	svc := newTestService()
	rs, _ := svc.CreateRuleSet(model.RuleSet{Name: "rs1"})
	cond, _ := svc.CreateCondition(model.Condition{Name: "c1", Field: "age", Operator: model.OpGE, Value: "18", ValueType: model.ValueTypeNumber})
	act, _ := svc.CreateAction(model.Action{Name: "a1", ActionType: model.ActionTypeLog})

	r, err := svc.CreateRule(model.Rule{
		RuleSetID:    rs.ID,
		Name:         "r1",
		Priority:     10,
		ConditionIDs: []string{cond.ID},
		ActionIDs:    []string{act.ID},
		Enabled:      true,
	})
	if err != nil {
		t.Fatalf("create rule: %v", err)
	}
	got, _ := svc.GetRule(r.ID)
	if got.Name != "r1" {
		t.Fatalf("name mismatch")
	}
	items, _, _ := svc.ListRules(model.RuleFilter{RuleSetID: rs.ID}, 1, 10)
	if len(items) != 1 {
		t.Fatalf("expected 1 rule")
	}
	_, _ = svc.ToggleRuleStatus(r.ID)
	got, _ = svc.GetRule(r.ID)
	if got.Status != model.RuleStatusInactive {
		t.Fatalf("expected inactive")
	}
}

func TestService_Evaluate(t *testing.T) {
	svc := newTestService()
	rs, _ := svc.CreateRuleSet(model.RuleSet{Name: "eval-rs"})
	cond, _ := svc.CreateCondition(model.Condition{Name: "age>=18", Field: "age", Operator: model.OpGE, Value: "18", ValueType: model.ValueTypeNumber})
	act, _ := svc.CreateAction(model.Action{Name: "log", ActionType: model.ActionTypeLog})
	_, _ = svc.CreateRule(model.Rule{
		RuleSetID:    rs.ID,
		Name:         "adult",
		Priority:     10,
		ConditionIDs: []string{cond.ID},
		ActionIDs:    []string{act.ID},
		Enabled:      true,
	})
	_, err := svc.EvaluateRuleSet(rs.ID, map[string]interface{}{"age": 20})
	if err == nil {
		t.Fatalf("expected error because ruleset is draft")
	}
	_, _ = svc.PublishRuleSet(rs.ID, "admin")
	er, err := svc.EvaluateRuleSet(rs.ID, map[string]interface{}{"age": 20})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if er.Result != model.ResultHit {
		t.Fatalf("expected hit")
	}
	er2, _ := svc.EvaluateRuleSet(rs.ID, map[string]interface{}{"age": 10})
	if er2.Result != model.ResultMiss {
		t.Fatalf("expected miss")
	}
}

func TestService_RuleVersionRollback(t *testing.T) {
	svc := newTestService()
	rs, _ := svc.CreateRuleSet(model.RuleSet{Name: "rollback"})
	_, _ = svc.PublishRuleSet(rs.ID, "admin")
	versions, _, _ := svc.ListRuleVersions(model.RuleVersionFilter{RuleSetID: rs.ID}, 1, 10)
	if len(versions) == 0 {
		t.Fatalf("expected version created")
	}
	_, err := svc.RollbackRuleVersion(versions[0].ID, "admin")
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}
	v, _ := svc.GetRuleVersion(versions[0].ID)
	if v.Status != model.RuleVersionStatusRolledBack {
		t.Fatalf("expected rolled_back")
	}
}

func TestService_ExportImport(t *testing.T) {
	svc := newTestService()
	rs, _ := svc.CreateRuleSet(model.RuleSet{Name: "export-me"})
	data, err := svc.ExportRuleSetSnapshot(rs.ID)
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	imported, err := svc.ImportRuleSetSnapshot(data, "admin")
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if imported.Name != "export-me-imported" {
		t.Fatalf("name mismatch")
	}
}

func TestService_Stats(t *testing.T) {
	svc := newTestService()
	_, _ = svc.CreateRuleSet(model.RuleSet{Name: "stats"})
	overview := svc.GetStatsOverview()
	if overview["rule_set_count"] != 1 {
		t.Fatalf("expected 1 ruleset")
	}
}

func TestService_ListPagination(t *testing.T) {
	svc := newTestService()
	for i := 0; i < 5; i++ {
		_, _ = svc.CreateRuleSet(model.RuleSet{Name: fmt.Sprintf("page-%d", i)})
	}
	items, total, _ := svc.ListRuleSets(model.RuleSetFilter{}, 1, 2)
	if len(items) != 2 {
		t.Fatalf("expected 2 items")
	}
	if total != 5 {
		t.Fatalf("expected total 5")
	}
}
