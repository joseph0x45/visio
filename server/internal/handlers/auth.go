package handlers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"log/slog"
	"net/http"
	"os"
	"time"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"
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

func (h *AuthHandler) GetAuthURL(c *fiber.Ctx) error {
	url := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		os.Getenv("GH_CLIENT_ID"),
		os.Getenv("GH_REDIRECT_URI"),
	)
	response := struct {
		URL string `json:"url"`
	}{
		URL: url,
	}
	if err := c.Status(fiber.StatusOK).JSON(response); err != nil {
		h.logger.Error(err.Error())
		return c.SendStatus(fiber.ErrInternalServerError.Code)
	}
	return nil
}

func (h *AuthHandler) GithubAuthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	error := c.Query("error")
	webAppURL := os.Getenv("WEB_APP_URL")
	errorURL := "%s/error?context=%s"
	internalErrRedirect := fmt.Sprintf(errorURL, webAppURL, "internal")
	if error != "" {
		var redirectURL string
		switch error {
		case "access_denied":
			redirectURL = fmt.Sprintf(errorURL, webAppURL, error)
		default:
			h.logger.Debug(fmt.Sprintf("Error while handling github redirect: %s", error))
			redirectURL = fmt.Sprint(errorURL, webAppURL, "unknown")
		}
		if err := c.Redirect(redirectURL, fiber.StatusTemporaryRedirect); err != nil {
			h.logger.Error(err.Error())
			return c.SendStatus(fiber.ErrInternalServerError.Code)
		}
		return nil
	}
	accessToken, err := pkg.GetToken(code)
	return nil
}

func (h *AuthHandler) GetUserInfo(c *fiber.Ctx) error {
	currentUser, ok := c.Context().Value("currentUser").(map[string]string)
	if !ok {
		return c.SendStatus(fiber.ErrUnauthorized.Code)
	}
	if err := c.Status(fiber.StatusOK).JSON(currentUser); err != nil {
		h.logger.Error(err.Error())
		return c.SendStatus(fiber.ErrInternalServerError.Code)
	}
	return nil
}
