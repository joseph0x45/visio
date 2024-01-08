package handlers

import (
	"fmt"
	"log/slog"
	"visio/internal/store"
	"visio/internal/types"

	"github.com/gofiber/fiber/v2"
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

func (h *AppHandler) GetLandingPage(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func (h *AppHandler) GetAuthPage(c *fiber.Ctx) error {
	return c.Render("auth", fiber.Map{})
}

func (h *AppHandler) GetHomePage(c *fiber.Ctx) error {
	return c.Render("home", fiber.Map{})
}

func (h *AppHandler) GetKeysPage(c *fiber.Ctx) error {
	currentUser, ok := c.Locals("currentUser").(*types.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	userKeys, err := h.keys.GetByUserId(currentUser.Id)
	if err != nil {
		h.logger.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
  fmt.Printf("%#v", userKeys)
	return c.Render("keys", fiber.Map{"Keys": userKeys}, "layouts/app")
}
