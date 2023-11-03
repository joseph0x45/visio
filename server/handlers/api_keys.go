package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"visio/repositories"
)

type KeyHandler struct {
	keys_repo  *repositories.Keys_repo
	logger     *logrus.Logger
	token_auth *jwtauth.JWTAuth
}

func NewKeyHandler(logger *logrus.Logger, keys_repo *repositories.Keys_repo, token_auth *jwtauth.JWTAuth) *KeyHandler {
	return &KeyHandler{
		keys_repo:  keys_repo,
		logger:     logger,
		token_auth: token_auth,
	}
}

func (h *KeyHandler) CreateKey(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	user_id := claims["user_id"].(string)
	if user_id == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

}

func (h *KeyHandler) RevokeKey(w http.ResponseWriter, r *http.Request) {

}

func (h *KeyHandler) RegisterRoutes(r chi.Router) {
  r.Get("/keys", func(w http.ResponseWriter, r *http.Request) {
    current_user := r.Context().Value("current_user").(map[string]string)
    println(current_user["id"])
    return
  })
}
