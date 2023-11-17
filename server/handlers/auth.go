package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"visio/models"
	"visio/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	logger             *logrus.Logger
	user_repo          *repositories.UserRepo
	githubOauth_config *oauth2.Config
	token_auth         *jwtauth.JWTAuth
}

func NewAuthHandler(logger *logrus.Logger, repo *repositories.UserRepo, githubOauthConfig *oauth2.Config, token_auth *jwtauth.JWTAuth) *AuthHandler {
	return &AuthHandler{
		logger:             logger,
		user_repo:          repo,
		githubOauth_config: githubOauthConfig,
		token_auth:         token_auth,
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
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, fmt.Sprintf("%s/error", os.Getenv("CONSOLE_URL")), http.StatusTemporaryRedirect)
		return
	}
	github_token, err := h.githubOauth_config.Exchange(r.Context(), code)
	if err != nil {
		h.logger.Error(err)
		http.Redirect(w, r, fmt.Sprintf("%s/error", os.Getenv("CONSOLE_URL")), http.StatusTemporaryRedirect)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		h.logger.Error(err)
		http.Redirect(w, r, fmt.Sprintf("%s/error", os.Getenv("CONSOLE_URL")), http.StatusTemporaryRedirect)
		return
	}
	req.Header.Add("Authorization", "Bearer "+github_token.AccessToken)
	response, err := client.Do(req)
	if err != nil {
		h.logger.Error(err)
		http.Redirect(w, r, fmt.Sprintf("%s/error", os.Getenv("CONSOLE_URL")), http.StatusTemporaryRedirect)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		h.logger.Error(err)
		http.Redirect(w, r, fmt.Sprintf("%s/error", os.Getenv("CONSOLE_URL")), http.StatusTemporaryRedirect)
		return
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		h.logger.Error(err)
		http.Redirect(w, r, fmt.Sprintf("%s/error", os.Getenv("CONSOLE_URL")), http.StatusTemporaryRedirect)
		return
	}
	username, avatar, github_id := data["name"].(string), data["avatar_url"].(string), data["id"].(float64)
	existing_user, err := h.user_repo.GetByGithubId(fmt.Sprintf("%.f", github_id))
	if err != nil {
		if err == sql.ErrNoRows {
			new_user_id := uuid.NewString()
			var new_user = &models.User{
				Id:       new_user_id,
				GithubId: fmt.Sprintf("%.f", github_id),
				Username: username,
				Avatar:   avatar,
			}
			err = h.user_repo.InsertNewUser(new_user)
			if err != nil {
				h.logger.Error(err)
				http.Redirect(w, r, fmt.Sprintf("%s/error", os.Getenv("CONSOLE_URL")), http.StatusTemporaryRedirect)
				return
			}
			_, auth_token, err := h.token_auth.Encode(
				map[string]interface{}{
					"user_id": new_user_id,
				},
			)
			if err != nil {
				h.logger.Error(err)
				http.Redirect(w, r, fmt.Sprintf("%s/error", os.Getenv("CONSOLE_URL")), http.StatusTemporaryRedirect)
				return
			}
			redirection_url := fmt.Sprintf("%s/login?token=%s", os.Getenv("CONSOLE_URL"), auth_token)
			http.Redirect(w, r, redirection_url, http.StatusTemporaryRedirect)
			return
		}
		h.logger.Error(err)
		http.Redirect(w, r, fmt.Sprintf("%s/error", os.Getenv("CONSOLE_URL")), http.StatusTemporaryRedirect)
		return
	}
	err = h.user_repo.UpdateUserInfos(fmt.Sprintf("%.f", github_id), username, avatar)
	if err != nil {
		h.logger.Error(err)
		http.Redirect(w, r, fmt.Sprintf("%s/error", os.Getenv("CONSOLE_URL")), http.StatusTemporaryRedirect)
		return
	}
	_, auth_token, err := h.token_auth.Encode(
		map[string]interface{}{
			"user_id": existing_user.Id,
		},
	)
	if err != nil {
		h.logger.Error(err)
		http.Redirect(w, r, fmt.Sprintf("%s/error", os.Getenv("CONSOLE_URL")), http.StatusTemporaryRedirect)
		return
	}
	redirection_url := fmt.Sprintf("%s/login?token=%s", os.Getenv("CONSOLE_URL"), auth_token)
	http.Redirect(w, r, redirection_url, http.StatusTemporaryRedirect)
	return
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/request", h.RequestGithubAuth)
	r.Get("/callback", h.GithubAuth)
}
