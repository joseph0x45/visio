package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"visio/internal/database"
	"visio/internal/handlers"
	"visio/internal/middlewares"
	"visio/internal/store"

	"github.com/Kagami/go-face"
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
	redisClient := database.GetRedisClient()
	users := store.NewUsersStore(postgresPool)
	sessions := store.NewSessionsStore(redisClient)
	keys := store.NewKeysStore(postgresPool)
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	appLogger := slog.New(textHandler)
	appHandler := handlers.NewAppHandler(keys, appLogger)
	authHandler := handlers.NewAuthHandler(users, sessions, appLogger)
	keyHandler := handlers.NewKeyHandler(keys, sessions, appLogger)
	recognizer, err := face.NewRecognizer(os.Getenv("MODELS_DIR"))
	if err != nil {
		panic(fmt.Sprintf("Error while initializing recognizer: %s", err.Error()))
	}
	faceHandler := handlers.NewFaceHandler(appLogger, recognizer)
	authMiddleware := middlewares.NewAuthMiddleware(sessions, users, appLogger)
	uploadMiddleware := middlewares.NewUploadMiddleware(appLogger)

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
	r.Get("/auth", appHandler.RenderAuthPage)

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth", authHandler.Authenticate)
	})

	r.Route("/keys", func(r chi.Router) {
		r.With(authMiddleware.CookieAuth).Get("/", appHandler.GetKeysPage)
		r.With(authMiddleware.CookieAuth).Post("/", keyHandler.Create)
		r.With(authMiddleware.CookieAuth).Delete("/{prefix}", keyHandler.Revoke)
	})

	r.Route("/faces", func(r chi.Router) {
		r.With(uploadMiddleware.HandleUploads(1)).Post("/", faceHandler.SaveFace)
	})

	port := os.Getenv("PORT")
	if port == "" {
		panic("Unable to read PORT environment variable")
	}
	fmt.Printf("Server listening on port %s\n", port)
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
