package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
	"visio/internal/store"
	"visio/internal/types"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

type Middleware struct {
	logger    *slog.Logger
	users     *store.Users
	tokenAuth *jwtauth.JWTAuth
}

func NewMiddlewareService(logger *slog.Logger, users *store.Users, tokenAuth *jwtauth.JWTAuth) *Middleware {
	return &Middleware{
		logger:    logger,
		users:     users,
		tokenAuth: tokenAuth,
	}
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := parts[1]
		tokenData, err := m.tokenAuth.Decode(token)
		if err != nil {
			m.logger.Error(fmt.Sprintf("Error while decoding auth token: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userId := tokenData.Subject()
		if userId == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		dbUser, err := m.users.GetById(userId)
		if err != nil {
			if errors.Is(err, types.ErrUserNotFound) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			m.logger.Error(fmt.Sprintf("Error while fetching user from database: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userData := map[string]string{
			"id":       dbUser.Id,
			"email":    dbUser.Email,
			"avatar":   dbUser.Avatar,
			"username": dbUser.Username,
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
