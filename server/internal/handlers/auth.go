package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oklog/ulid/v2"
	"log/slog"
	"net/http"
	"os"
	"time"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"
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
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		if err != http.ErrNoCookie {
			h.logger.Debug(fmt.Sprintf("Error while reading cookies from client: %v", err))
		}
		goto returnAuthURL
	} else {
		sessionId := sessionCookie.Value
		sessionUser, err := h.sessions.Get(sessionId)
		if err != nil {
			if !errors.Is(err, types.ErrSessionNotFound) {
				h.logger.Debug(err.Error())
			}
			goto returnAuthURL
		}
		_, err = h.users.GetById(sessionUser)
		if err != nil {
			if !errors.Is(err, types.ErrUserNotFound) {
				h.logger.Error(err.Error())
			}
			goto returnAuthURL
		}
		http.Redirect(w, r, fmt.Sprintf("%s/dashboard", os.Getenv("WEB_APP_URL")), http.StatusPermanentRedirect)
		return
	}
returnAuthURL:
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
	webAppURL := os.Getenv("WEB_APP_URL")
	errorURL := "%s/error?context=%s"
	internalErrRedirect := fmt.Sprintf(errorURL, webAppURL, "internal")
	if error != "" {
		var redirectURL string
		switch error {
		case "access_denied":
			redirectURL = fmt.Sprintf(errorURL, webAppURL, error)
		default:
			h.logger.Debug(fmt.Sprintf("Error while handling github redirect: %s", error))
			redirectURL = fmt.Sprint(errorURL, webAppURL, "unknown")
		}
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		return
	}
	accessToken, err := pkg.GetToken(code)
	if err != nil {
		h.logger.Error(err.Error())
		redirectURL := fmt.Sprintf(errorURL, webAppURL, "internal")
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		return
	}
	userData, err := pkg.GetUserData(accessToken)
	if userData.Email == "" {
		userData.Email, err = pkg.GetUserPrimaryEmail(accessToken)
		if err != nil {
			if errors.Is(err, types.ErrNoPrimaryEmailFound) {
				redirectURL := fmt.Sprintf(errorURL, webAppURL, "no_mail_found")
				http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				return
			}
			h.logger.Error(err.Error())
			http.Redirect(w, r, internalErrRedirect, http.StatusTemporaryRedirect)
			return
		}
	}
	if err != nil {
		h.logger.Error(err.Error())
		http.Redirect(w, r, internalErrRedirect, http.StatusTemporaryRedirect)
		return
	}
	dbUser, err := h.users.GetByGithubId(fmt.Sprintf("%.f", userData.Id))
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			newUser := &types.User{
				Id:         ulid.Make().String(),
				GithubId:   fmt.Sprintf("%.f", userData.Id),
				Email:      userData.Email,
				Username:   userData.Login,
				Avatar:     userData.Avatar,
				SignupDate: time.Now().UTC(),
			}
			err = h.users.Insert(newUser)
			if err != nil {
				h.logger.Error(err.Error())
				http.Redirect(w, r, internalErrRedirect, http.StatusTemporaryRedirect)
				return
			}
			sessionId := ulid.Make().String()
			err = h.sessions.Create(sessionId, newUser.Id)
			if err != nil {
				h.logger.Error(err.Error())
				http.Redirect(w, r, internalErrRedirect, http.StatusTemporaryRedirect)
				return
			}
			loginRedirect := fmt.Sprintf("%s/login?session=%s", webAppURL, sessionId)
			http.Redirect(w, r, loginRedirect, http.StatusTemporaryRedirect)
			return
		}
		h.logger.Error(err.Error())
		http.Redirect(w, r, internalErrRedirect, http.StatusTemporaryRedirect)
		return
	}
	err = h.users.UpdateUserData(dbUser.Id, userData.Email, userData.Login, userData.Avatar)
	if err != nil {
		h.logger.Error(err.Error())
		http.Redirect(w, r, internalErrRedirect, http.StatusTemporaryRedirect)
		return
	}
	sessionId := ulid.Make().String()
	err = h.sessions.Create(sessionId, dbUser.Id)
	if err != nil {
		h.logger.Error(err.Error())
		http.Redirect(w, r, internalErrRedirect, http.StatusTemporaryRedirect)
		return
	}
	loginRedirect := fmt.Sprintf("%s/login?session=%s", webAppURL, sessionId)
	http.Redirect(w, r, loginRedirect, http.StatusTemporaryRedirect)
	return
}

func (h *AuthHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := r.Context().Value("currentUser").(map[string]interface{})
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userData, err := json.Marshal(currentUser)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userData)
	return
}
