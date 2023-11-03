package pkg

import "golang.org/x/crypto/bcrypt"

func Hash(text string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), 7)
	return string(hash), err
}

func HashMatches(hash string, text string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
	return err == nil
}
