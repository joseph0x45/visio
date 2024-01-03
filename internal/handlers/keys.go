package handlers

import (
	"log/slog"
	"visio/internal/store"

	"github.com/gofiber/fiber/v2"
)

type KeyHandler struct {
	users    *store.Users
	logger   *slog.Logger
	sessions *store.Sessions
}

func NewKeyHandler(usersStore *store.Users, sessionsStore *store.Sessions, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		users:    usersStore,
		sessions: sessionsStore,
		logger:   logger,
	}
}

func (h *KeyHandler) CreateKey(c *fiber.Ctx) error {
	// Check how many number of keys the user has already created

	// Generate prefix - ULID format

	// Generate 23 characters long string - the key

	// Hash the key

	// Store key and prefix in DB

	//Return key in format <prefix>.<key (unhased)> to user

	return nil
}
