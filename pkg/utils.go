package pkg

import (
  "os"
  "fmt"
  "math/rand"
)

func CleanupFiles(files []string) []error {
	var errors []error
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			errors = append(errors, fmt.Errorf("Error while deleting file %s: %w", file, err))
		}
	}
	return errors
}


func GenerateRandomString(length int) string {
	const CHARACTER_POOL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890-_=+:;?><|"
	key := ""
	for i := 0; i < length; i++ {
		idx := rand.Intn(len(CHARACTER_POOL))
		key += string(CHARACTER_POOL[idx])
	}
	return key
}

