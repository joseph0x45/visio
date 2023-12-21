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

func HashMatches(plainTextPassword, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainTextPassword))
	if err != nil {
		return fmt.Errorf("Error while comparing hash and password: %w", err)
	}
	return nil
}
