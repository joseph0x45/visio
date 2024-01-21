package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"visio/internal/store"
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
		"views/landing.html",
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

func (h *AppHandler) RenderHomePage(w http.ResponseWriter, r *http.Request){
	templateFiles := []string{
		"views/layouts/base.html",
		"views/home.html",
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
