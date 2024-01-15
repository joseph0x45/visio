package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

const (
	KeyLimit = 5
)

type KeyHandler struct {
	keys     *store.Keys
	logger   *slog.Logger
	sessions *store.Sessions
}

func NewKeyHandler(keysStore *store.Keys, sessionsStore *store.Sessions, logger *slog.Logger) *KeyHandler {
	return &KeyHandler{
		keys:     keysStore,
		sessions: sessionsStore,
		logger:   logger,
	}
}

func generateRandomString(length int) string {
	const CHARACTER_POOL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890-_=+:;'/?><|"
	key := ""
	for i := 0; i < length; i++ {
		idx := rand.Intn(len(CHARACTER_POOL))
		key += string(CHARACTER_POOL[idx])
	}
	return key
}

func (h *KeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := r.Context().Value("currentUser").(*types.User)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	keysCount, err := h.keys.CountByOwnerId(currentUser.Id)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if keysCount == KeyLimit {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	prefix := generateRandomString(7)
	suffix := generateRandomString(23)
	key := fmt.Sprintf("%s.%s", prefix, suffix)
	keyHash, err := pkg.Hash(suffix)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newKey := &types.Key{
		Id:           ulid.Make().String(),
		UserId:       currentUser.Id,
		Prefix:       prefix,
		KeyHash:      keyHash,
		CreationDate: time.Now().UTC().Format("January, 2 2006"),
	}
	err = h.keys.Insert(newKey)
	if err != nil {
		if errors.Is(err, types.ErrDuplicatePrefix) {
			h.logger.Debug("A duplicate prefix error occured")
		}
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(
		map[string]interface{}{
			"data": map[string]interface{}{
				"naked_key":  key,
				"key_object": newKey,
			},
		},
	)
	if err != nil {
		//Delete created key and add header to tell client
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
	return
}

func (h *KeyHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := r.Context().Value("currentUser").(*types.User)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	prefix := chi.URLParam(r, "prefix")
	if prefix == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.keys.Delete(prefix, currentUser.Id)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
