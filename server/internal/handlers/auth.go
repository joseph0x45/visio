package handlers

import (
	"net/http"
	"visio/internal/store"
)

type AuthHandler struct {
	users *store.Users
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) LoginWithGithub(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	error := r.URL.Query().Get("error")
	if error != "" {
		//user denied authorization
	}
	if code == "" {
		//bad request
	}

}
