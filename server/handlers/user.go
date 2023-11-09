package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	// "visio/models"
	"visio/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	logger    *logrus.Logger
	user_repo *repositories.UserRepo
}

func NewUserHandler(logger *logrus.Logger, user_repo *repositories.UserRepo) *UserHandler {
	return &UserHandler{
		logger:    logger,
		user_repo: user_repo,
	}
}

func (h *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	current_user, ok := r.Context().Value("current_user").(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := h.user_repo.GetById(current_user["id"])
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(
		map[string]string{
			"id":       user.Id,
			"avatar":   user.Avatar,
			"username": user.Username,
			"plan":     user.Plan,
			"email":    user.Email,
		},
	)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
  w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Get("/user", h.GetUserInfo)
}
