package middlewares

import (
	"context"
	"errors"
	"fmt"
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
        fmt.Println("No cookie? ;(")
				http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
				return
			}
		}
		sessionValue, err := m.sessions.Get(sessionCookie.Value)
		if err != nil {
			if errors.Is(err, types.ErrSessionNotFound) {
        fmt.Println("No session 0_0")
				http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
				return
			}
			m.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sessionUser, err := m.users.GetById(sessionValue)
		if err != nil {
			if errors.Is(err, types.ErrUserNotFound) {
        fmt.Println("No user o_o")
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
