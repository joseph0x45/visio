package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"visio/internal/store"
	"visio/internal/types"
)

type AuthMiddleware struct {
	sessions *store.Sessions
	users    *store.Users
	logger   *slog.Logger
}

func NewAuthMiddleware(sessions *store.Sessions, users *store.Users, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		sessions: sessions,
		users:    users,
		logger:   logger,
	}
}

func (m *AuthMiddleware) CookieAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		sessionValue, err := m.sessions.Get(sessionCookie.Name)
		if err != nil {
			if errors.Is(err, types.ErrSessionNotFound) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			m.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sessionUser, err := m.users.GetById(sessionValue)
		if err != nil {
			if errors.Is(err, types.ErrUserNotFound) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			m.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "currentUser", sessionUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
