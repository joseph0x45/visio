package handlers

import (
	"api/pkg"
	"api/stores"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	users  *stores.Users
	logger *logrus.Logger
}

func NewAuthHandler(u *stores.Users, l *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		users:  u,
		logger: l,
	}
}

func (h *AuthHandler) RequestAuth(w http.ResponseWriter, r *http.Request) {
	oauthURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email",
		os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_OAUTH_REDIRECT_URI"),
	)
	data, err := json.Marshal(
		map[string]string{
			"url": oauthURL,
		},
	)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}

func (h *AuthHandler) HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	error := r.URL.Query().Get("error")
	if error != "" {
		//Something went wrong
		if error == "access_denied" {
			//User denied access
		}
	}
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	accessToken, err := pkg.ExchangeCodeWithToken(code, h.logger)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := pkg.GetUserGithubData(accessToken, h.logger)
  _ = data
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/auth/request", h.RequestAuth)
	r.Get("/auth/github/callback", h.HandleAuthCallback)
}
