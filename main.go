package main

import (
	"fmt"
	"log/slog"
	"os"
	"visio/internal/database"
	"visio/internal/handlers"
	"visio/internal/middlewares"
	"visio/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	appEnv := os.Getenv("ENV")
	if appEnv != "PROD" {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}
	appHandler := handlers.NewAppHandler()
	postgresPool := database.NewPostgresPool()
	redisClient := database.GetRedisClient()
	users := store.NewUsersStore(postgresPool)
	sessions := store.NewSessionsStore(redisClient)
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	appLogger := slog.New(textHandler)
	authHandler := handlers.NewAuthHandler(users, sessions, appLogger)
	authMiddleware := middlewares.NewAuthMiddleware(sessions, users, appLogger)

	engine := html.New("./views", ".html")
	engine.Reload(appEnv != "PROD")
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})
	app.Static("/public", "./public")
	app.Use(recover.New())

	client := app.Group("/")
	client.Get("/", appHandler.GetLandingPage)
	client.Get("/auth", appHandler.GetAuthPage)
	client.Get("/home", authMiddleware.CookieAuth, appHandler.GetHomePage)

	server := app.Group("/api")
	server.Post("/auth", authHandler.Signup)

	port := os.Getenv("PORT")
	if port == "" {
		panic("Unable to read PORT environment variable")
	}
	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
