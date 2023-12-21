package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"visio/internal/database"
	"visio/internal/handlers"
	"visio/internal/middlewares"
	"visio/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func (s *Server) RegisterRoutes() http.Handler {
	pgPool := database.NewPostgresPool()
	usersStore := store.NewUsersStore(pgPool)
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(jsonHandler)
	loggingMiddleware := middlewares.NewLoggingMiddleware(logger)
  tokenAuth:= jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	authHandler := handlers.NewAuthHandler(usersStore, logger, tokenAuth)

	r := chi.NewRouter()
	r.Use(loggingMiddleware.SpamFilter)
	r.Use(loggingMiddleware.RequestLogger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello from health")
		w.WriteHeader(http.StatusOK)
		return
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})
	return r
}
