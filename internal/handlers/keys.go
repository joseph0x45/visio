package handlers

import (
	"log/slog"
	"math/rand"
	"net/http"
	"visio/internal/store"
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

}
