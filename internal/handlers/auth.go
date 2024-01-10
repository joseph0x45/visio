package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"
	"github.com/oklog/ulid/v2"
)

type AuthHandler struct {
	users    *store.Users
	logger   *slog.Logger
	sessions *store.Sessions
}

func NewAuthHandler(usersStore *store.Users, sessionsStore *store.Sessions, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		users:    usersStore,
		sessions: sessionsStore,
		logger:   logger,
	}
}

func (h *AuthHandler) Authicate(w http.ResponseWriter, r *http.Request) {
	payload := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	action := r.URL.Query().Get("action")
  fmt.Println(action)
	if action != "Login" && action != "Register" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch action {
	case "Register":
		count, err := h.users.CountByEmail(payload.Email)
		if err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if count > 0 {
			w.WriteHeader(http.StatusConflict)
			return
		}
		hash, err := pkg.Hash(payload.Password)
		if err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		newUser := &types.User{
			Id:           ulid.Make().String(),
			Email:        payload.Email,
			PasswordHash: hash,
			SignupDate:   time.Now().UTC().Format("January, 2 2006"),
		}
		err = h.users.Insert(newUser)
		if err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sessionId := ulid.Make().String()
		err = h.sessions.Create(sessionId, newUser.Id)
		if err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		authCookie := &http.Cookie{
			Name:  "session",
			Value: sessionId,
		}
		http.SetCookie(w, authCookie)
		w.WriteHeader(http.StatusCreated)
		return
	case "Login":
		dbUser, err := h.users.GetByEmail(payload.Email)
		if err != nil {
			if errors.Is(err, types.ErrUserNotFound) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !pkg.PasswordMatches(payload.Password, dbUser.PasswordHash) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sessionId := ulid.Make().String()
		err = h.sessions.Create(sessionId, dbUser.Id)
		if err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		authCookie := &http.Cookie{
			Name:  "session",
			Value: sessionId,
		}
		http.SetCookie(w, authCookie)
		w.WriteHeader(http.StatusOK)
		return
	}
}
