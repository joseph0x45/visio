package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"log/slog"
	"net/http"
	"os"
	"visio/internal/database"
	"visio/internal/handlers"
	"visio/internal/middlewares"
	"visio/internal/store"
)

func (s *Server) RegisterRoutes() http.Handler {
	pgPool := database.NewPostgresPool()
	redisClient := database.GetRedisClient()
	usersStore := store.NewUsersStore(pgPool)
	sessionsStore := store.NewSessionsStore(redisClient)
	jsonHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(jsonHandler)
	loggingMiddleware := middlewares.NewLoggingMiddleware(logger)

	authHandler := handlers.NewAuthHandler(usersStore, sessionsStore, logger)

	r := chi.NewRouter()
	// r.Use(loggingMiddleware.SpamFilter)
	r.Use(loggingMiddleware.RequestLogger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello from health")
		w.WriteHeader(http.StatusOK)
		return
	})

	r.Route("/auth", func(r chi.Router) {
		r.Get("/callback", authHandler.GithubAuthCallback)
	})
	return r
}
