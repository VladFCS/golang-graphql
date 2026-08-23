package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/vladfc/ghira/internal/config"
)

func New(cfg config.HTTPConfig, logger *slog.Logger, graphqlHandler http.Handler) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler)

	mux.Handle("POST /query", graphqlHandler)
	handler := recoverer(logger)(requestLogger(logger)(mux))

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()

			next.ServeHTTP(w, r)

			logger.Info(
				"http request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(startedAt).String(),
			)
		})
	}
}

func recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if value := recover(); value != nil {
					logger.Error("panic recovered", "value", value)
					writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("failed to write json response", "error", err)
	}
}
