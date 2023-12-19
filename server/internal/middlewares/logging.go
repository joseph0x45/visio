package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type LoggerMiddleware struct {
	logger *slog.Logger
}

func NewLoggingMiddleware(logger *slog.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
    logger: logger,
  }
}

func (m *LoggerMiddleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.logger.Info("Started processing request 1")
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		m.logger.Info(fmt.Sprintf("Request 1 took %s", elapsed.String()))
    return
	})
}
