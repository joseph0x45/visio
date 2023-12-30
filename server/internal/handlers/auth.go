package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type AuthHandler struct {
	users     *store.Users
	logger    *slog.Logger
}

func NewAuthHandler(usersStore *store.Users, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		users:     usersStore,
		logger:    logger,
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

func (h *AuthHandler) GithubAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	error := r.URL.Query().Get("error")
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
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		return
	}
	accessToken, err := pkg.GetToken(code)
	if err != nil {
		h.logger.Error(err.Error())
		redirectURL := fmt.Sprintf(errorURL, webAppURL, "internal")
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		return
	}
	userData, err := pkg.GetUserData(accessToken)
	if userData.Email == "" {
		userData.Email, err = pkg.GetUserPrimaryEmail(accessToken)
		if err != nil {
			if errors.Is(err, types.ErrNoPrimaryEmailFound) {
				redirectURL := fmt.Sprintf(errorURL, webAppURL, "no_mail_found")
				http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				return
			}
			h.logger.Error(err.Error())
			http.Redirect(w, r, internalErrRedirect, http.StatusTemporaryRedirect)
			return
		}
	}
	if err != nil {
		h.logger.Error(err.Error())
		http.Redirect(w, r, internalErrRedirect, http.StatusTemporaryRedirect)
		return
	}
	dbUser, err := h.users.GetByGithubId(fmt.Sprintf("%.f", userData.Id))
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			newUser := &types.User{
				Id:         ulid.Make().String(),
				GithubId:   fmt.Sprintf("%.f", userData.Id),
				Email:      userData.Email,
				Username:   userData.Login,
				Avatar:     userData.Avatar,
				SignupDate: time.Now().UTC(),
			}
			err = h.users.Insert(newUser)
			if err != nil {
				h.logger.Error(err.Error())
				http.Redirect(w, r, internalErrRedirect, http.StatusTemporaryRedirect)
				return
			}
			_, authToken, err := h.tokenAuth.Encode(map[string]interface{}{
				"sub": newUser.Id,
			})
			if err != nil {
				h.logger.Error(fmt.Sprintf("Error while encoding auth token: %v", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			loginRedirect := fmt.Sprintf("%s/login?token=%s", webAppURL, authToken)
			http.Redirect(w, r, loginRedirect, http.StatusTemporaryRedirect)
			return
		}
		h.logger.Error(err.Error())
		http.Redirect(w, r, internalErrRedirect, http.StatusTemporaryRedirect)
		return
	}
	err = h.users.UpdateUserData(dbUser.Id, userData.Email, userData.Login, userData.Avatar)
	if err != nil {
		h.logger.Error(err.Error())
		http.Redirect(w, r, internalErrRedirect, http.StatusTemporaryRedirect)
		return
	}
	_, authToken, err := h.tokenAuth.Encode(map[string]interface{}{
		"sub": dbUser.Id,
	})
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while encoding auth token: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	loginRedirect := fmt.Sprintf("%s/login?token=%s", webAppURL, authToken)
	http.Redirect(w, r, loginRedirect, http.StatusTemporaryRedirect)
	return
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
