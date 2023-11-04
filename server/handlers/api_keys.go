package handlers

import (
	"encoding/json"
	"net/http"
	"visio/models"
	"visio/pkg"
	"visio/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
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
  w.Header().Add("Content-Type", "application/json")
	w.Write(data)
	return
}

func (h *KeyHandler) CreateKey(w http.ResponseWriter, r *http.Request) {
	current_user, ok := r.Context().Value("current_user").(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user_keys_total, err := h.keys_repo.GetUserNumberOfKeys(current_user["id"])
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if current_user["plan"] == "basic" {
		if user_keys_total == 1 {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	key_prefix := pkg.GenerateRandomString(7)
	key_suffix := pkg.GenerateRandomString(23)
	key_hash, err := pkg.Hash(key_suffix)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	new_key := &models.Key{
		Id:      uuid.NewString(),
		Owner:   current_user["id"],
		Prefix:  key_prefix,
		KeyHash: key_hash,
	}
	err = h.keys_repo.InsertNewKey(new_key)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(
		map[string]string{
			"key": key_prefix + "." + key_suffix,
		},
	)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
	return
}

func (h *KeyHandler) RevokeKey(w http.ResponseWriter, r *http.Request) {

}

func (h *KeyHandler) RegisterRoutes(r chi.Router) {
	r.Get("/keys", h.GetKeys)
  r.Post("/keys", h.CreateKey)
}
