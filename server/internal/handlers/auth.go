package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"

	"github.com/oklog/ulid/v2"
)

type AuthHandler struct {
	users    *store.Users
	sessions *store.Sessions
	logger   *slog.Logger
}

func NewAuthHandler(usersStore *store.Users, sessionsStore *store.Sessions, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		users:    usersStore,
		sessions: sessionsStore,
		logger:   logger,
	}
}

func (h *AuthHandler) GetAuthURL(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		os.Getenv("GH_CLIENT_ID"),
		os.Getenv("GH_REDIRECT_URI"),
	)
	data, err := json.Marshal(
		map[string]string{
			"url": url,
		},
	)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

func (h *AuthHandler) GithubAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	error := r.URL.Query().Get("error")
	if error != "" {
		if error == "access_denied" {
			http.Redirect(w, r, "http://localhost:5173/error?context=access_denied", http.StatusTemporaryRedirect)
			return
		}
		h.logger.Debug(fmt.Sprintf("Error while handling github redirect %s", error))
		http.Redirect(w, r, "http://localhost:5173/error?context=unknown", http.StatusTemporaryRedirect)
		return
	}
	accessToken, err := pkg.GetToken(code)
	if err != nil {
		h.logger.Error(err.Error())
		http.Redirect(w, r, "http://localhost:5173/error?context=internal", http.StatusTemporaryRedirect)
		return
	}
	userData, err := pkg.GetUserData(accessToken)
	if err != nil {
		h.logger.Error(err.Error())
		http.Redirect(w, r, "http://localhost:5173/error?context=internal", http.StatusTemporaryRedirect)
		return
	}
	dbUser, err := h.users.GetByGithubId(fmt.Sprintf("%.f", userData.Id))
	if err != nil {
		if errors.Is(err, types.ErrNoUserFound) {
			newUser := &types.User{
				Id:         ulid.Make().String(),
				GithubId:   fmt.Sprintf("%.f", userData.Id),
				Username:   userData.Login,
				Avatar:     userData.Avatar,
				SignupDate: time.Now().UTC(),
			}
			err = h.users.Insert(newUser)
			if err != nil {
				h.logger.Error(err.Error())
				http.Redirect(w, r, "http://localhost:5173/error?context=internal", http.StatusTemporaryRedirect)
				return
			}
			sessionId := ulid.Make().String()
			err = h.sessions.Create(sessionId, newUser.Id)
			if err != nil {
				h.logger.Error(err.Error())
				http.Redirect(w, r, "http://localhost:5173/error?context=internal", http.StatusTemporaryRedirect)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("http://localhost:5173/login?session=%s", sessionId), http.StatusTemporaryRedirect)
			return
		}
		h.logger.Error(err.Error())
		http.Redirect(w, r, "http://localhost:5173/error?context=internal", http.StatusTemporaryRedirect)
		return
	}
	err = h.users.UpdateUserData(dbUser.Id, userData.Login, userData.Avatar)
	if err != nil {
		h.logger.Error(err.Error())
		http.Redirect(w, r, "http://localhost:5173/error?context=internal", http.StatusTemporaryRedirect)
		return
	}
	sessionId := ulid.Make().String()
	err = h.sessions.Create(sessionId, dbUser.Id)
	if err != nil {
		h.logger.Error(err.Error())
		http.Redirect(w, r, "http://localhost:5173/error?context=internal", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("http://localhost:5173/login?session=%s", sessionId), http.StatusTemporaryRedirect)
	return
}
