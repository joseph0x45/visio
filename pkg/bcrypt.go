package pkg

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hash(plainTextPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error while hashing password: %w", err)
	}
	return string(hash), nil
}

func PasswordMatches(plainTextPassword, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainTextPassword)) == nil
}
