package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"time"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type AuthHandler struct {
	users    *store.Users
	logger   *slog.Logger
	sessions *store.Sessions
}

func NewAuthHandler(usersStore *store.Users, sessionsStore *store.Sessions, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		users:    usersStore,
		sessions: sessionsStore,
		logger:   logger,
	}
}

func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	reqPayload := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	action := c.Query("action")
	if action == "" {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}
	if err := c.BodyParser(reqPayload); err != nil {
		h.logger.Error(fmt.Sprintf("Error while parsing body: %v", err))
		return c.SendStatus(fiber.ErrInternalServerError.Code)
	}
	switch action {
	case "Register":
		count, err := h.users.CountByEmail(reqPayload.Email)
		if err != nil {
			h.logger.Error(err.Error())
			return c.SendStatus(fiber.ErrInternalServerError.Code)
		}
		if count != 0 {
			return c.SendStatus(fiber.ErrConflict.Code)
		}
		hash, err := pkg.Hash(reqPayload.Password)
		if err != nil {
			h.logger.Error(err.Error())
			return c.SendStatus(fiber.ErrInternalServerError.Code)
		}
		newUser := &types.User{
			Id:         ulid.Make().String(),
			Email:      reqPayload.Email,
			Password:   hash,
			SignupDate: time.Now().UTC(),
		}
		err = h.users.Insert(newUser)
		if err != nil {
			h.logger.Error(err.Error())
			return c.SendStatus(fiber.ErrInternalServerError.Code)
		}
		sessionId := ulid.Make().String()
		err = h.sessions.Create(sessionId, newUser.Id)
		if err != nil {
			h.logger.Error(err.Error())
			c.Set("X-Err-Context", "ERR_SESSION_CREATION")
			return c.SendStatus(fiber.ErrInternalServerError.Code)
		}
		authCookie := &fiber.Cookie{
			Name:  "session",
			Value: sessionId,
		}
		c.Cookie(authCookie)
		return c.SendStatus(fiber.StatusCreated)
	case "Login":
		dbUser, err := h.users.GetByEmail(reqPayload.Email)
		if err != nil {
			if errors.Is(err, types.ErrUserNotFound) {
				return c.SendStatus(fiber.ErrBadRequest.Code)
			}
			h.logger.Error(err.Error())
			return c.SendStatus(fiber.ErrInternalServerError.Code)
		}
		if !pkg.PasswordMatches(reqPayload.Password, dbUser.Password) {
			fmt.Print("bad pwd\n")
			return c.SendStatus(fiber.ErrBadRequest.Code)
		}
		sessionId := ulid.Make().String()
		err = h.sessions.Create(sessionId, dbUser.Id)
		if err != nil {
			h.logger.Error(err.Error())
			c.Set("X-Err-Context", "ERR_SESSION_CREATION")
			return c.SendStatus(fiber.ErrInternalServerError.Code)
		}
		authCookie := &fiber.Cookie{
			Name:  "session",
			Value: sessionId,
		}
		c.Cookie(authCookie)
		return c.SendStatus(fiber.StatusOK)
	}
	return c.SendStatus(fiber.ErrBadRequest.Code)
}
