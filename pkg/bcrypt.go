package pkg

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hash(plainTextString string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextString), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error while hashing string: %w", err)
	}
	return string(hash), nil
}

func HashMatches(plainTextString, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainTextString)) == nil
}
