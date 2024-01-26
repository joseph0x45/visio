package store

import (
	"os"
	"testing"
	"visio/internal/database"
)

func SessionsMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestSession(t *testing.T) {
	sessionManager := database.NewSessionManager()
	t.Run("Create session", func(t *testing.T) {
		sessionManager.CreateSession("session", "userId")
	})

	t.Run("Get session", func(t *testing.T) {
		sessionValue := sessionManager.GetSession("session")
		if sessionValue != "userId" {
			t.Fatalf("Wanted %s got %s", "userId", sessionValue)
		}
	})

	t.Run("Get non existing session", func(t *testing.T) {
		sessionValue := sessionManager.GetSession("random")
		if sessionValue != "" {
			t.Fatalf("Wanted empty string got %s", sessionValue)
		}
	})

	t.Run("Delete session", func(t *testing.T) {
		sessionManager.DeleteSession("session")
		v := sessionManager.GetSession("session")
		if v != "" {
			t.Fatalf("Wanted empty string got %s", v)
		}
	})
}
