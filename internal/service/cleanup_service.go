package service

import (
	"time"

	"ruleengine/internal/model"
)

func (s *Service) CleanupOldEvaluationRecords(before time.Time) (int, error) {
	all := s.store.ListEvaluationRecords()
	deleted := 0
	for _, er := range all {
		if er.EvaluatedAt.Before(before) {
			if err := s.store.DeleteEvaluationRecord(er.ID); err == nil {
				deleted++
			}
		}
	}
	return deleted, nil
}

func (s *Service) CleanupOldExecutionLogs(before time.Time) (int, error) {
	all := s.store.ListExecutionLogs()
	deleted := 0
	for _, el := range all {
		if el.ExecutedAt.Before(before) {
			if err := s.store.DeleteExecutionLog(el.ID); err == nil {
				deleted++
			}
		}
	}
	return deleted, nil
}

func (s *Service) CleanupOldAuditRecords(before time.Time) (int, error) {
	all := s.store.ListAuditRecords()
	deleted := 0
	for _, ar := range all {
		if ar.OperatedAt.Before(before) {
			if err := s.store.DeleteAuditRecord(ar.ID); err == nil {
				deleted++
			}
		}
	}
	return deleted, nil
}

func (s *Service) CleanupOldRuleVersions(keepLatest int) (int, error) {
	if keepLatest < 1 {
		keepLatest = 5
	}
	all := s.store.ListRuleVersions()
	grouped := make(map[string][]*model.RuleVersion)
	for _, rv := range all {
		grouped[rv.RuleSetID] = append(grouped[rv.RuleSetID], rv)
	}
	deleted := 0
	for _, list := range grouped {
		if len(list) <= keepLatest {
			continue
		}
		for i := 0; i < len(list)-keepLatest; i++ {
			if list[i].Status == model.RuleVersionStatusArchived {
				if err := s.store.DeleteRuleVersion(list[i].ID); err == nil {
					deleted++
				}
			}
		}
	}
	return deleted, nil
}
