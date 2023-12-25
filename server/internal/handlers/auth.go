package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	users     *store.Users
	logger    *slog.Logger
	tokenAuth *jwtauth.JWTAuth
}

func NewAuthHandler(usersStore *store.Users, logger *slog.Logger, tokenAuth *jwtauth.JWTAuth) *AuthHandler {
	return &AuthHandler{
		users:  usersStore,
		logger: logger,
    tokenAuth: tokenAuth,
	}
}

func (h *AuthHandler) GithubAuthCallback(w http.ResponseWriter, r *http.Request){

}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value("requestId").(string)
	payload := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while decoding body %v", err), "requestId", requestId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	count, err := h.users.CountByEmail(payload.Email)
	if err != nil {
		h.logger.Error(err.Error(), "requestId", r.Context().Value("requestId"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if count > 0 {
		w.WriteHeader(http.StatusConflict)
		return
	}
	hash, err := pkg.Hash(payload.Password)
	if err != nil {
		h.logger.Error(err.Error(), "requestId", requestId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUser := &types.User{
		Id:         uuid.NewString(),
		Email:      payload.Email,
		Password:   hash,
		SignupDate: time.Now(),
	}
	err = h.users.Insert(newUser)
	if err != nil {
		h.logger.Error(err.Error(), "requestId", requestId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, token, err := h.tokenAuth.Encode(map[string]interface{}{
		"userId": newUser.Id,
	})
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while encoding JWT: %v", err), "requestId", requestId)
    w.Header().Set("X-Failure-Reason", "Token-Generation-Failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while marshalling response: %v", err), "requestId", requestId)
    w.Header().Set("X-Failure-Reason", "Data-Marshalling-Failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
	return
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

}
