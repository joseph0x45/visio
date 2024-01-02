package handlers

import (
	"fmt"
	"log/slog"
	"visio/internal/store"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	users  *store.Users
	logger *slog.Logger
}

func NewAuthHandler(usersStore *store.Users, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		users:  usersStore,
		logger: logger,
	}
}

func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	reqPayload := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		h.logger.Error(fmt.Sprintf("Error while parsing body: %v", err))
		return c.SendStatus(fiber.ErrInternalServerError.Code)
	}
	// count, err := h.users.CountByEmail(reqPayload.Email)
	// if err != nil {
	// 	h.logger.Error(err.Error())
	// 	return c.SendStatus(fiber.ErrInternalServerError.Code)
	// }
	// if count != 0 {
	// 	return c.SendStatus(fiber.ErrConflict.Code)
	// }
	return nil
}
