package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"os"
)

type (
	Host struct {
		Fiber *fiber.App
	}
)

func main() {
	appEnv := os.Getenv("ENV")
	if appEnv != "PROD" {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}
	hosts := map[string]*Host{}

	api := fiber.New()
	api.Use(logger.New())
	api.Use(recover.New())

	hosts["api.127.0.0.1:8080"] = &Host{api}

	engine := html.New("./views", ".html")
	engine.Reload(appEnv != "PROD")
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})
	app.Static("/public", "./public")
	app.Use(logger.New())
	app.Use(recover.New())

	hosts["127.0.0.1:8080"] = &Host{app}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Hello": "Bozo",
		})
	})

	server := fiber.New()
	server.Use(func(c *fiber.Ctx) error {
		host := hosts[c.Hostname()]
		fmt.Println(c.Hostname())
		if host == nil {
			return c.SendStatus(fiber.ErrNotFound.Code)
		} else {
			host.Fiber.Handler()(c.Context())
			return nil
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		panic("Unable to read PORT environment variable")
	}
	err := server.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
