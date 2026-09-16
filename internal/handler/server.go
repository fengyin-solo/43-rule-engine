// Package handler 实现 HTTP 处理器层。
package handler

import (
	"errors"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"ruleengine/internal/config"
	"ruleengine/internal/model"
	"ruleengine/internal/service"
	"ruleengine/internal/store"
	"ruleengine/pkg/httpx"
	"ruleengine/pkg/logger"
)

type Server struct {
	svc *service.Service
	log *logger.Logger
	cfg *config.Config
}

func NewServer(svc *service.Service, log *logger.Logger, cfg *config.Config) *Server {
	return &Server{svc: svc, log: log, cfg: cfg}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	s.registerRuleSetRoutes(mux)
	s.registerRuleRoutes(mux)
	s.registerConditionRoutes(mux)
	s.registerActionRoutes(mux)
	s.registerEvaluationRecordRoutes(mux)
	s.registerRuleVersionRoutes(mux)
	s.registerExecutionLogRoutes(mux)
	s.registerDataTypeRoutes(mux)
	s.registerAuditRecordRoutes(mux)
	s.registerRuleStatRoutes(mux)
	s.registerEvalRoutes(mux)
	s.registerStatsRoutes(mux)
	s.registerExportImportRoutes(mux)
	s.registerBatchRoutes(mux)
	s.registerReportRoutes(mux)
	s.registerHealthRoutes(mux)
	s.registerCleanupRoutes(mux)
	mux.Handle("GET /", http.FileServer(http.Dir("web")))
	return s.rateLimitMiddleware(s.authMiddleware(s.loggingMiddleware(s.recoveryMiddleware(mux))))
}

func (s *Server) maxPageSize() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 100
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.log.Errorf("panic: %v\n%s", rec, debug.Stack())
				httpx.InternalError(w, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" || r.URL.Path == "/style.css" || r.URL.Path == "/app.js" {
			next.ServeHTTP(w, r)
			return
		}
		key := r.Header.Get("X-Api-Key")
		if key == "" {
			key = r.URL.Query().Get("api_key")
		}
		if key != s.cfg.APIKey {
			httpx.Unauthorized(w, "无效的 API Key")
			return
		}
		next.ServeHTTP(w, r)
	})
}

type tokenBucket struct {
	tokens   float64
	capacity float64
	rate     float64
	last     time.Time
	mu       sync.Mutex
}

func newTokenBucket(capacity, rate float64) *tokenBucket {
	return &tokenBucket{tokens: capacity, capacity: capacity, rate: rate, last: time.Now()}
}

func (tb *tokenBucket) allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(tb.last).Seconds()
	tb.tokens += elapsed * tb.rate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.last = now
	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	tb := newTokenBucket(100, 10)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !tb.allow() {
			httpx.Error(w, http.StatusTooManyRequests, 429, "请求过于频繁")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case model.IsValidationError(err):
		httpx.BadRequest(w, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.NotFound(w, err.Error())
	case errors.Is(err, store.ErrConflict):
		httpx.Conflict(w, err.Error())
	default:
		httpx.InternalError(w, err.Error())
	}
}
