package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"visio/internal/database"
	"visio/internal/handlers"
	// "visio/internal/middlewares"
	"visio/internal/store"
)

func (s *Server) RegisterRoutes() http.Handler {
	pgPool := database.NewPostgresPool()
	usersStore := store.NewUsersStore(pgPool)
	jsonHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(jsonHandler)
	// loggingMiddleware := middlewares.NewLoggingMiddleware(logger)

	authHandler := handlers.NewAuthHandler(usersStore, logger)

	r := chi.NewRouter()
	// r.Use(loggingMiddleware.SpamFilter)
	// r.Use(loggingMiddleware.RequestLogger)
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
