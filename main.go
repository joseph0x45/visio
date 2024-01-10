package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"visio/internal/database"
	"visio/internal/handlers"
	"visio/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

//go:embed views/*
var views embed.FS

//go:embed public/output.css
var publicFS embed.FS

func main() {
	appEnv := os.Getenv("ENV")
	if appEnv != "PROD" {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}
	postgresPool := database.NewPostgresPool()
	// redisClient := database.GetRedisClient()
	// users := store.NewUsersStore(postgresPool)
	// sessions := store.NewSessionsStore(redisClient)
	keys := store.NewKeysStore(postgresPool)
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	appLogger := slog.New(textHandler)
	appHandler := handlers.NewAppHandler(keys, appLogger)
	// authHandler := handlers.NewAuthHandler(users, sessions, appLogger)
	// keyHandler := handlers.NewKeyHandler(keys, sessions, appLogger)
	// authMiddleware := middlewares.NewAuthMiddleware(sessions, users, appLogger)

	r := chi.NewRouter()

	r.Get("/public/output.css", func(w http.ResponseWriter, r *http.Request) {
		css, err := publicFS.ReadFile("public/output.css")
		if err != nil {
			fmt.Printf("Error while reading css file from embedded file system: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(""))
			return
		}
		w.Header().Set("Content-Type", "text/css")
		w.Write(css)
		return
	})

	r.Get("/", appHandler.RenderLandingPage)

	// engine := html.New("./views", ".html")
	// engine.Reload(appEnv != "PROD")
	// engine.AddFunc("jsonify", func(s interface{}) string {
	// 	jsonBytes, err := json.Marshal(s)
	// 	if err != nil {
	// 		return ""
	// 	}
	// 	return string(jsonBytes)
	// })
	// app := fiber.New(fiber.Config{
	// 	Views:       engine,
	// 	ViewsLayout: "layouts/main",
	// })
	// app.Static("/public", "./public")
	// app.Use(recover.New())
	//
	// client := app.Group("/")
	// client.Get("/", appHandler.GetLandingPage)
	// client.Get("/auth", appHandler.GetAuthPage)
	// client.Get("/keys", authMiddleware.CookieAuth, appHandler.GetKeysPage)
	//
	// server := app.Group("/api")
	// server.Post("/auth", authHandler.Signup)
	// server.Post("/key", authMiddleware.CookieAuth, keyHandler.Create)
	// server.Delete("/key/:prefix", authMiddleware.CookieAuth, keyHandler.Revoke)

	port := os.Getenv("PORT")
	if port == "" {
		panic("Unable to read PORT environment variable")
	}
	// err := app.Listen(fmt.Sprintf(":%s", port))
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
