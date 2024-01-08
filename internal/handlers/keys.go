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

func generateRandomString(length int) string {
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
	if err != nil {
		h.logger.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if key_count > KEY_LIMIT {
		return c.SendStatus(fiber.StatusForbidden)
	}
	prefix := generateRandomString(7)
	suffix := generateRandomString(23)
	hashedKey, err := pkg.Hash(suffix)
	if err != nil {
		h.logger.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	generatedKey := &types.Key{
		Id:           ulid.Make().String(),
		UserId:       currentUser.Id,
		Prefix:       prefix,
		KeyHash:      hashedKey,
		CreationDate: time.Now().UTC().Format("January, 2 2006"),
	}
	if err := h.keys.Insert(generatedKey); err != nil {
		if errors.Is(err, types.ErrDuplicatePrefix) {
			h.logger.Debug("Duplicate prefix error triggered")
		}
		h.logger.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	key := fmt.Sprintf("%s.%s", prefix, suffix)
	err = c.JSON(
		map[string]interface{}{
			"data": map[string]string{
				"key": key,
			},
		},
	)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusCreated)
}
