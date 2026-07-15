package httpserver

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/example/open-company-brain/apps/brain-api/internal/config"
	"github.com/example/open-company-brain/apps/brain-api/internal/health"
)

type Server struct {
	config  config.RuntimeConfig
	checker *health.Checker
	logger  *slog.Logger
}

func New(cfg config.RuntimeConfig, logger *slog.Logger) *Server {
	return &Server{
		config:  cfg,
		logger:  logger,
		checker: health.New(cfg.DependencyTimeout, cfg.PostgresAddress, cfg.RedisAddress, cfg.RedisPassword, cfg.QdrantURL, cfg.MinIOHealthURL, cfg.WorkerHealthURL),
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /live", s.live)
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /api/v1/instance", s.instance)
	return loggingMiddleware(s.logger, mux)
}

func (s *Server) live(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "alive"})
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	result := s.checker.Check(r.Context(), s.config.Brain.Instance.Name)
	code := http.StatusOK
	if result.Status != "healthy" {
		code = http.StatusServiceUnavailable
	}
	writeJSON(w, code, result)
}

func (s *Server) instance(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"name":     s.config.Brain.Instance.Name,
		"language": s.config.Brain.Instance.Language,
		"modules":  s.config.Brain.Modules,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("http_request", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(started).Milliseconds())
	})
}
