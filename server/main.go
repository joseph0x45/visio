package main

import (
	"fmt"
	"net/http"
	"os"
	"visio/handlers"
	"visio/pkg"
	"visio/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	godotenv.Load()
	logger := logrus.New()
	logger.SetReportCaller(true)
	db, err := sqlx.Connect("postgres", os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
  r := chi.NewRouter()
  r.Use(middleware.Logger)
  r.Use(cors.Handler(cors.Options{
    AllowedOrigins:   []string{"https://*", "http://*"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    ExposedHeaders:   []string{"Link"},
    AllowCredentials: false,
    MaxAge:           300,
  }))

	githubOauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	users_repo := repositories.NewUserRepo(db)
  keys_repo := repositories.NewKeysRepo(db)

	auth_handler := handlers.NewAuthHandler(logger, users_repo, githubOauthConfig, tokenAuth)
  keys_handler := handlers.NewKeyHandler(logger, keys_repo, tokenAuth)

  middleware_service := pkg.NewAuthMiddlewareService(tokenAuth, users_repo)

  r.Route("/auth", func(r chi.Router) {
    auth_handler.RegisterRoutes(r)
  })

  // Authenticated routes
  r.Route("/", func(r chi.Router) {
    r.Use(middleware_service.Authenticate)
    keys_handler.RegisterRoutes(r)
  })


	fmt.Println("Server launched on port 8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
