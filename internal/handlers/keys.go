package handlers

import (
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
	// Check how many number of keys the user has already created

	// Generate prefix - ULID format
	prefix := ulid.Make().String()

	// Generate 23 characters long string - the key
	key := generateKey(23)

	// Hash the key
	hashedKey, err := pkg.Hash(key)
	if err != nil {
		h.logger.Error(err.Error())
		return c.SendStatus(fiber.ErrInternalServerError.Code)
	}

	currentUser := c.Locals("currentUser").(*types.User)
	generatedKey := &types.Key{
		Owner:           currentUser.Id,
		Prefix:          prefix,
		Key:             key,
		KeyCreationDate: time.Now().UTC(),
	}

	_ = hashedKey
	fmt.Printf("%+v", generatedKey)

	// Store key and prefix in DB

	//Return key in format <prefix>.<key (unhased)> to user

	finalKey := fmt.Sprintf("%s.%s", prefix, hashedKey)
	return c.Send([]byte(finalKey))
}
