package handlers

import "github.com/gofiber/fiber/v2"

type AppHandler struct {
}

func NewAppHandler() *AppHandler {
	return &AppHandler{}
}

func (h *AppHandler) GetLandingPage(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func (h *AppHandler) GetAuthPage(c *fiber.Ctx) error {
	return c.Render("auth", fiber.Map{})
}
