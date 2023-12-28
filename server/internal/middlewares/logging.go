package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	"visio/internal/store"
	"visio/internal/types"

	"github.com/google/uuid"
)

type Middleware struct {
	logger   *slog.Logger
	users    *store.Users
	sessions *store.Sessions
}

func NewMiddlewareService(logger *slog.Logger, users *store.Users, sessions *store.Sessions) *Middleware {
	return &Middleware{
		logger:   logger,
		users:    users,
		sessions: sessions,
	}
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			m.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sessionUser, err := m.sessions.Get(sessionCookie.Value)
		if err != nil {
			if errors.Is(err, types.ErrSessionNotFound) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			m.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		dbUser, err := m.users.GetById(sessionUser)
		if err != nil {
			if errors.Is(err, types.ErrUserNotFound) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			m.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userData := map[string]string{
			"id":       dbUser.Id,
			"email":    dbUser.Email,
			"username": dbUser.Username,
			"avatar":   dbUser.Avatar,
		}
		ctx := context.WithValue(r.Context(), "currentUser", userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
