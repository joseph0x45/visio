package pkg

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"visio/repositories"

	"github.com/go-chi/jwtauth/v5"
)

type MiddlewareService struct {
	token_auth *jwtauth.JWTAuth
	users_repo *repositories.UserRepo
}

func NewAuthMiddlewareService(token_auth *jwtauth.JWTAuth, users_repo *repositories.UserRepo) *MiddlewareService {
	return &MiddlewareService{
		token_auth: token_auth,
		users_repo: users_repo,
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
