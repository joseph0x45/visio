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
	"github.com/Kagami/go-face"
)

func main() {
	godotenv.Load()
	rec, err := face.NewRecognizer(os.Getenv("MODELS_DIR"))
	if err!= nil {
		panic(err)
	}
	defer rec.Close()
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
		RedirectURL:  "https://api.getvisio.cloud/auth/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	users_repo := repositories.NewUserRepo(db)
	keys_repo := repositories.NewKeysRepo(db)
	faces_repo := repositories.NewFacesRepo(db)

	auth_handler := handlers.NewAuthHandler(logger, users_repo, githubOauthConfig, tokenAuth)
  user_handler := handlers.NewUserHandler(logger, users_repo)
	keys_handler := handlers.NewKeyHandler(logger, keys_repo, tokenAuth)
	faces_handler_v1 := handlers.NewFacesHandlerV1(logger, faces_repo, rec)

	middleware_service := pkg.NewAuthMiddlewareService(tokenAuth, users_repo, keys_repo, logger)

	r.Route("/auth", func(r chi.Router) {
		auth_handler.RegisterRoutes(r)
	})

	r.Route("/", func(r chi.Router) {
		r.Use(middleware_service.Authenticate)
		keys_handler.RegisterRoutes(r)
    user_handler.RegisterRoutes(r)
	})

  r.Route("/v1", func(r chi.Router) {
    r.Use(middleware_service.AuthenticateWithKey)
    faces_handler_v1.RegisterRoutes(r)
  })

	fmt.Println("Server launched on port 8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
