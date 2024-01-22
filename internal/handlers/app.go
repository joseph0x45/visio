package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"

	"github.com/oklog/ulid/v2"
)

type AppHandler struct {
	users  *store.Users
	keys   *store.Keys
	logger *slog.Logger
}

func NewAppHandler(users *store.Users, keys *store.Keys, logger *slog.Logger) *AppHandler {
	return &AppHandler{
		users:  users,
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

func (h *AppHandler) RenderHomePage(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := r.Context().Value("currentUser").(*types.User)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userKey, err := h.keys.GetByUserId(currentUser.Id)
	nakedKey := ""
	userHasKey := true
	if err != nil {
		if !errors.Is(err, types.ErrKeyNotFound) {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userHasKey = false
	}
	if !userHasKey {
		prefix := pkg.GenerateRandomString(7)
		suffix := pkg.GenerateRandomString(23)
		nakedKey = fmt.Sprintf("%s.%s", prefix, suffix)
		keyHash, err := pkg.Hash(suffix)
		if err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userKey = &types.Key{
			Id:           ulid.Make().String(),
			UserId:       currentUser.Id,
			Prefix:       prefix,
			KeyHash:      keyHash,
			CreationDate: time.Now().UTC().Format("January, 2 2006"),
		}
		err = h.keys.Insert(userKey)
		if err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	data := map[string]any{
		"NewKey":   !userHasKey,
		"NakedKey": nakedKey,
		"Prefix":   userKey.Prefix,
	}
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
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *AppHandler) RevokeKey(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := r.Context().Value("currentUser").(*types.User)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err := h.keys.Delete(currentUser.Id)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
