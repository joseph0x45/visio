package pkg

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
)

func ToJPEG(filePath string) (string, error) {
	imgBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("Error while reading file: %w", err)
	}
	pngBytes, err := png.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return "", fmt.Errorf("Error while decoding PNG bytes: %w", err)
	}
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, pngBytes, nil); err != nil {
		return "", fmt.Errorf("Error while encoding bytes: %w", err)
	}
	f, err := os.CreateTemp("", "")
	if err != nil {
		return "", fmt.Errorf("Error while creating temp file: %w", err)
	}
	err = os.WriteFile(f.Name(), buffer.Bytes(), os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("Error while writing bytes %w", err)
	}
	return f.Name(), nil
}
