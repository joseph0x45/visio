package middlewares

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type Middleware struct {
	logger *slog.Logger
}

func NewLoggingMiddleware(logger *slog.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (m *Middleware) SpamFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appItendifier := r.Header.Get("X-VISIO-APP-IDENTIFIER")
		authHeader := r.Header.Get("Authorization")
		if appItendifier != os.Getenv("IDENTIFIER") {
			if authHeader == "" {
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqId := uuid.NewString()
		m.logger.Info(fmt.Sprintf("Started processing request %s", reqId))
		ctx := context.WithValue(r.Context(), "requestId", reqId)
		next.ServeHTTP(w, r.WithContext(ctx))
		elapsed := time.Since(start)
		m.logger.Info(fmt.Sprintf("Request %s took %s", reqId, elapsed.String()))
		return
	})
}
