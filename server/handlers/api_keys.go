package handlers

import (
	"encoding/json"
	"net/http"
	"visio/repositories"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sirupsen/logrus"
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

func (h *KeyHandler) GetKeys(w http.ResponseWriter, r *http.Request) {
	current_user, ok := r.Context().Value("current_user").(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	keys, err := h.keys_repo.SelectKeys(current_user["id"])
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(keys)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
  w.Write(data)
  return
}

func (h *KeyHandler) CreateKey(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *KeyHandler) RevokeKey(w http.ResponseWriter, r *http.Request) {

}

func (h *KeyHandler) RegisterRoutes(r chi.Router) {
	r.Get("/keys", h.GetKeys)
}
