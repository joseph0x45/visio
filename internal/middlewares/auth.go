package middlewares

import (
	"errors"
	"log/slog"
	"visio/internal/store"
	"visio/internal/types"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	sessions *store.Sessions
	users    *store.Users
	logger   *slog.Logger
}

func NewAuthMiddleware(sessions *store.Sessions, users *store.Users, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		sessions: sessions,
		users:    users,
		logger:   logger,
	}
}

func (m *AuthMiddleware) CookieAuth(c *fiber.Ctx) error {
	sessionId := c.Cookies("session")
	if sessionId == "" {
		return c.Redirect("/auth", fiber.StatusFound)
	}
	sessionValue, err := m.sessions.Get(sessionId)
	if err != nil {
		if errors.Is(err, types.ErrSessionNotFound) {
			return c.Redirect("/auth", fiber.StatusFound)
		}
		m.logger.Error(err.Error())
		return c.Redirect("/auth", fiber.StatusFound)
	}
	sessionUser, err := m.users.GetById(sessionValue)
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			return c.Redirect("/auth", fiber.StatusFound)
		}
		m.logger.Error(err.Error())
		return c.Redirect("/auth", fiber.StatusFound)
	}

	c.Locals("currentUser", sessionUser)
	return c.Next()
}

func (m *AuthMiddleware) KeyAuth(c *fiber.Ctx) error {
	return c.Next()
}
