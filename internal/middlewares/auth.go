package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"
)

type AuthMiddleware struct {
	sessions *store.Sessions
	users    *store.Users
	keys     *store.Keys
	logger   *slog.Logger
}

func NewAuthMiddleware(sessions *store.Sessions, users *store.Users, keys *store.Keys, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		sessions: sessions,
		users:    users,
		keys:     keys,
		logger:   logger,
	}
}

func (m *AuthMiddleware) CookieAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
				return
			}
		}
		sessionValue := m.sessions.Get(sessionCookie.Value)
		sessionUser, err := m.users.GetById(sessionValue)
		if err != nil {
			if errors.Is(err, types.ErrUserNotFound) {
				http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
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

func (m *AuthMiddleware) KeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		parts := strings.Split(apiKey, ".")
		if len(parts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		prefix, suffix := parts[0], parts[1]
		key, err := m.keys.GetByPrefix(prefix)
		if err != nil {
			if errors.Is(err, types.ErrKeyNotFound) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			m.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		hashMatches := pkg.HashMatches(suffix, key.KeyHash)
		if !hashMatches {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "currentUser", key.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
