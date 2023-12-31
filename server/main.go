package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	hosts["api.localhost:8080"] = &Host{api}

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	hosts["localhost:8080"] = &Host{app}

	app.Get("/home", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	server := fiber.New()
	server.Use(func(c *fiber.Ctx) error {
		host := hosts[c.Hostname()]
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
	fmt.Printf("Server listening on port %s", port)
	err := server.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
