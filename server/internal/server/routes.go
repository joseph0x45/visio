package server

import (
	"log/slog"
	"net/http"
	"os"
	"visio/internal/database"
	"visio/internal/handlers"
	"visio/internal/middlewares"
	"visio/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
)

func (s *Server) RegisterRoutes() http.Handler {
	pgPool := database.NewPostgresPool()
	usersStore := store.NewUsersStore(pgPool)
	jwtAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
	jsonHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(jsonHandler)
	middlewareService := middlewares.NewMiddlewareService(logger, usersStore, jwtAuth)
	authHandler := handlers.NewAuthHandler(usersStore, jwtAuth, logger)

	r := chi.NewRouter()
	r.Use(middlewareService.RequestLogger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-VISIO-APP-IDENTIFIER"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
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
