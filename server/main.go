package main

import (
	"fmt"
	"net/http"
	"os"
	"visio/handlers"
	"visio/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	githubOauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

  tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

  users_repo := repositories.NewUserRepo(db)

  auth_handler := handlers.NewAuthHandler(logger, users_repo, githubOauthConfig, tokenAuth)

  auth_handler.RegisterRoutes(r)

  fmt.Println("Server launched on port 8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
