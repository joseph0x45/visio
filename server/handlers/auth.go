package handlers

import (
	"encoding/json"
	"net/http"
	"visio/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	logger    *logrus.Logger
	user_repo *repositories.UserRepo
  githubOauth_config *oauth2.Config
}

func NewAuthHandler(logger *logrus.Logger, repo *repositories.UserRepo, githubOauthConfig *oauth2.Config) *AuthHandler {
	return &AuthHandler{
		logger:    logger,
		user_repo: repo,
    githubOauth_config: githubOauthConfig,
	}
}

func (h *AuthHandler) RequestGithubAuth(w http.ResponseWriter, r *http.Request) {
	url := h.githubOauth_config.AuthCodeURL("state", oauth2.SetAuthURLParam("client_id", h.githubOauth_config.ClientID))
	data, err := json.Marshal(
		map[string]string{
			"url": url,
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

func (h *AuthHandler) GithubAuth(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) RegisterRoutes(r *chi.Mux) {
	r.Get("/auth/request", h.RequestGithubAuth)
	r.Get("/auth/callback", h.GithubAuth)
	r.Post("/register", h.Register)
}
