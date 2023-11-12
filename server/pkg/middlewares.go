package pkg

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"visio/repositories"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sirupsen/logrus"
)

type MiddlewareService struct {
	token_auth *jwtauth.JWTAuth
	logger     *logrus.Logger
	users_repo *repositories.UserRepo
	keys_repo  *repositories.Keys_repo
}

func NewAuthMiddlewareService(token_auth *jwtauth.JWTAuth, users_repo *repositories.UserRepo, keys_repo *repositories.Keys_repo, logger *logrus.Logger) *MiddlewareService {
	return &MiddlewareService{
		token_auth: token_auth,
		users_repo: users_repo,
		keys_repo:  keys_repo,
		logger:     logger,
	}
}

func (m *MiddlewareService) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth_header := r.Header.Get("Authorization")
		if len(strings.Split(auth_header, " ")) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		bearer_token := strings.Split(auth_header, " ")[1]
		claims, err := m.token_auth.Decode(bearer_token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user_id, ok := claims.Get("user_id")
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user, err := m.users_repo.GetById(user_id.(string))
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			m.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "current_user", map[string]string{
			"id":   user.Id,
			"plan": user.Plan,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *MiddlewareService) AuthenticateWithKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api_key := r.Header.Get("Authorization")
		if len(strings.Split(api_key, ".")) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		key_prefix := strings.Split(api_key, ".")[0]
		key, err := m.keys_repo.GetKeyByPrefix(key_prefix)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			m.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		current_user, err := m.users_repo.GetById(key.Owner)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			m.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "current_user", map[string]string{
			"id":   current_user.Id,
			"plan": current_user.Plan,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
