package main

import (
	"embed"
	"fmt"
	"github.com/Kagami/go-face"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"visio/internal/database"
	"visio/internal/handlers"
	"visio/internal/middlewares"
	"visio/internal/store"
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
	users := store.NewUsersStore(postgresPool)
	sessionManager := database.NewSessionManager()
	sessions := store.NewSessionsStore(sessionManager)
	keys := store.NewKeysStore(postgresPool)
	faces := store.NewFacesStore(postgresPool)
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	appLogger := slog.New(textHandler)
	appHandler := handlers.NewAppHandler(users, keys, appLogger)
	authHandler := handlers.NewAuthHandler(users, sessions, appLogger)
	recognizer, err := face.NewRecognizer(os.Getenv("MODELS_DIR"))
	if err != nil {
		panic(fmt.Sprintf("Error while initializing recognizer: %s", err.Error()))
	}
	faceHandler := handlers.NewFaceHandler(appLogger, recognizer, faces)
	authMiddleware := middlewares.NewAuthMiddleware(sessions, users, keys, appLogger)
	uploadMiddleware := middlewares.NewUploadMiddleware(appLogger)

	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)

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
	r.With(authMiddleware.CookieAuth).Get("/home", appHandler.RenderHomePage)

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth", authHandler.Authenticate)
	})

	r.Route("/faces", func(r chi.Router) {
		r.With(authMiddleware.KeyAuth).With(uploadMiddleware.HandleUploads(1)).Post("/", faceHandler.SaveFace)
		r.With(authMiddleware.KeyAuth).Delete("/{id}", faceHandler.DeleteFace)
		r.With(authMiddleware.KeyAuth).Get("/", faceHandler.GetAll)
		r.With(authMiddleware.KeyAuth).Get("/{id}", faceHandler.GetById)
		r.With(authMiddleware.KeyAuth).Route("/compare", func(r chi.Router) {
			r.With(uploadMiddleware.HandleUploads(2)).Post("/", faceHandler.CompareUploaded)
			r.With(uploadMiddleware.HandleUploads(0)).Post("/saved", faceHandler.CompareSavedFaces)
			r.With(uploadMiddleware.HandleUploads(1)).Post("/mixed", faceHandler.CompareMixt)
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		panic("Unable to read PORT environment variable")
	}
	fmt.Printf("Server listening on port %s\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		panic(err)
	}
}
