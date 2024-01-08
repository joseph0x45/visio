package handlers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"log/slog"
	"math/rand"
	"time"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"
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
	currentUser, ok := c.Locals("currentUser").(*types.User)
	if !ok {
		err := errors.New("Error during currentUser type conversion")
		h.logger.Error(err.Error())
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	const KEY_LIMIT = 3
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
	prefix := ulid.Make().String()
	key := generateKey(23)
	hashedKey, err := pkg.Hash(key)
	if err != nil {
		h.logger.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	generatedKey := &types.Key{
		UserId:       currentUser.Id,
		Prefix:       prefix,
		KeyHash:      hashedKey,
		CreationDate: time.Now().UTC(),
	}
	if err := h.keys.Insert(generatedKey); err != nil {
		h.logger.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	finalKey := fmt.Sprintf("%s.%s", prefix, key)
	err = c.JSON(
		map[string]interface{}{
			"data": map[string]string{
				"key": finalKey,
			},
		},
	)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusCreated)
}
