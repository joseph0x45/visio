package pkg

import (
  "os"
  "fmt"
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
