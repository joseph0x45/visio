package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
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

func (h *AppHandler) RenderLandingPage(w http.ResponseWriter, r *http.Request) {
	templFiles := []string{
		"views/layouts/base.html",
		"views/home.html",
	}
	ts, err := template.ParseFiles(templFiles...)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
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
	return c.Render("keys", fiber.Map{"Keys": userKeys}, "layouts/app")
}
