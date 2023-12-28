package server

import (
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
	middlewareService := middlewares.NewMiddlewareService(logger, usersStore, sessionsStore)
	authHandler := handlers.NewAuthHandler(usersStore, sessionsStore, logger)

	r := chi.NewRouter()
	r.Use(middlewareService.RequestLogger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/auth", func(r chi.Router) {
		r.Get("/callback", authHandler.GithubAuthCallback)
		r.Group(func(r chi.Router) {
			r.Use(middlewareService.SpamFilter)
			r.Get("/url", authHandler.GetAuthURL)
			r.With(middlewareService.Authenticate).Get("/user", authHandler.GetUserInfo)
		})
	})
	return r
}
