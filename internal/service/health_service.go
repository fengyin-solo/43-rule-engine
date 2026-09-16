package service

import (
	"time"
)

func (s *Service) HealthCheck() map[string]interface{} {
	return map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	}
}

func (s *Service) CheckDependencies() map[string]string {
	result := make(map[string]string)
	result["store"] = "ok"
	return result
}
