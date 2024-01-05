package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
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

func generateKey(length int) string {
	const CHARACTER_POOL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890-_=+:;'/?><|"
	key := ""

	for i := 0; i < length; i++ {
		idx := rand.Intn(len(CHARACTER_POOL))
		key += string(CHARACTER_POOL[idx])
	}

	return key
}

func (h *KeyHandler) CreateKey(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*types.User)
	const KEY_LIMIT = 3

	// Check how many number of keys the user has already created
	key_count, err := h.keys.CountByOwnerId(currentUser.Id)

	fmt.Println("Key count", key_count)
	if err != nil {
		h.logger.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if key_count > KEY_LIMIT {
		err := errors.New("Limit of keys exceeded")
		h.logger.Error(err.Error())
		return c.SendStatus(fiber.StatusForbidden)
	}

	// Generate prefix - ULID format
	prefix := ulid.Make().String()

	// Generate 23 characters long string - the key
	key := generateKey(23)

	// Hash the key
	hashedKey, err := pkg.Hash(key)
	if err != nil {
		h.logger.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	generatedKey := &types.Key{
		KeyOwner:        currentUser.Id,
		Prefix:          prefix,
		KeyHash:         hashedKey,
		KeyCreationDate: time.Now().UTC(),
	}

	_ = hashedKey

	// Store key and prefix in DB
	if err := h.keys.Insert(generatedKey); err != nil {
		h.logger.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	//Return key in format <prefix>.<key (unhased)> to user
	finalKey := fmt.Sprintf("%s.%s", prefix, key)
	return c.Send([]byte(finalKey))
}
