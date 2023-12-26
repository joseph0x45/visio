package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"
)

type AuthHandler struct {
	users  *store.Users
	logger *slog.Logger
}

func NewAuthHandler(usersStore *store.Users, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		users:  usersStore,
		logger: logger,
	}
}

func (h *AuthHandler) GithubAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	error := r.URL.Query().Get("error")
	if error != "" {
		if error == "access_denied" {
			//Handle user denied auth
		}
	}
	accessToken, err := pkg.GetToken(code)
	if err != nil {
		h.logger.Error(err.Error())
		println("err 1")
		http.Redirect(w, r, fmt.Sprintf("http://localhost:5173/error"), http.StatusTemporaryRedirect)
		return
	}
	userData, err := pkg.GetUserData(accessToken)
	if err != nil {
		h.logger.Error(err.Error())
		println("err 2")
		http.Redirect(w, r, fmt.Sprintf("http://localhost:5173/error"), http.StatusTemporaryRedirect)
		return
	}
	if userData.Email == "" {
		userData.Email, err = pkg.GetUserPrimaryEmail(accessToken)
		if err != nil {
			if errors.Is(err, types.ErrNoPrimaryEmailFound) {
				h.logger.Error("No primary email found for ser")
				println("err 3")
				http.Redirect(w, r, fmt.Sprintf("http://localhost:5173/error"), http.StatusTemporaryRedirect)
				return
			}
			h.logger.Error(err.Error())
			println("err 4")
			http.Redirect(w, r, fmt.Sprintf("http://localhost:5173/error"), http.StatusTemporaryRedirect)
			return
		}
	}
	sessionId := "yay"
	println("going there")
	http.Redirect(w, r, fmt.Sprintf("http://localhost:5173/login?session=%s", sessionId), http.StatusTemporaryRedirect)
	return
}
