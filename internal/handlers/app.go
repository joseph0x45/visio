package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"visio/internal/store"
	"visio/internal/types"
)

type AppHandler struct {
	keys   *store.Keys
	logger *slog.Logger
}

func NewAppHandler(keys *store.Keys, logger *slog.Logger) *AppHandler {
	return &AppHandler{
		keys:   keys,
		logger: logger,
	}
}

func (h *AppHandler) RenderLandingPage(w http.ResponseWriter, r *http.Request) {
	templFiles := []string{
		"views/layouts/base.html",
		"views/home.html",
	}
	ts, err := template.ParseFiles(templFiles...)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *AppHandler) RenderAuthPage(w http.ResponseWriter, r *http.Request) {
	templateFiles := []string{
		"views/layouts/base.html",
		"views/auth.html",
	}
	ts, err := template.ParseFiles(templateFiles...)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *AppHandler) GetKeysPage(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := r.Context().Value("currentUser").(*types.User)
	if !ok {
		http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
		return
	}
	userKeys, err := h.keys.GetByUserId(currentUser.Id)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	templateFiles := []string{
		"views/layouts/app.html",
		"views/keys.html",
	}
	ts, err := template.ParseFiles(templateFiles...)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	templData := map[string]interface{}{
		"Keys": userKeys,
	}
	err = ts.ExecuteTemplate(w, "app", templData)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
