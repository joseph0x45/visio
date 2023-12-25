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
  println("received request")
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
		return
	}
	userData, err := pkg.GetUserData(accessToken)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	test := userData.Email == ""
	h.logger.Info(fmt.Sprintf("%v", test))
  h.logger.Info(userData.Email)
	if userData.Email == "" {
		userData.Email, err = pkg.GetUserPrimaryEmail(accessToken)
		if err != nil {
			if errors.Is(err, types.ErrNoPrimaryEmailFound) {
				h.logger.Error("No primary email found for ser")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	return
}
