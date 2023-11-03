package handlers

import (
	"net/http"
	"visio/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
  logger *logrus.Logger
  user_repo *repositories.UserRepo
}

func NewAuthHandler(logger *logrus.Logger, repo *repositories.UserRepo) *AuthHandler {
  return &AuthHandler{
    logger: logger,
    user_repo: repo,
  }
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request){

}

func (h *AuthHandler) RegisterRoutes(r *chi.Mux){
  r.Post("/register", h.Register)
}
