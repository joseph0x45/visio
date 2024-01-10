package pkg

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"visio/internal/types"
)

var imageMimeTypes = []string{
	"image/jpeg",
	"image/png",
}

const MaxFileSize = 5

func HandleFileUpload(w http.ResponseWriter, r *http.Request) (string, error) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize<<20)
	reader, err := r.MultipartReader()
	if err != nil {
		return "", fmt.Errorf("Error while reading multipart request body: %w", err)
	}
	p, err := reader.NextPart()
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("Error while reading multipart body part: %w", err)
	}
	if p.FormName() != "face" {
		return "", types.ErrFileNotFound
	}
	buffer := bufio.NewReader(p)
	sniffedBytes, err := buffer.Peek(512)
	if err != nil {
		return "", fmt.Errorf("Error while peeking through bytes: %w", err)
	}
	contentType := http.DetectContentType(sniffedBytes)
	fileIsNotImage := true
	for _, mimeType := range imageMimeTypes {
		if contentType == mimeType {
			fileIsNotImage = false
			break
		}
	}
	if fileIsNotImage {
		return "", types.ErrUnsupportedFormat
	}
	f, err := os.CreateTemp("", "")
	if err != nil {
		return "", fmt.Errorf("Error while creating file: %w", err)
	}
	defer f.Close()
	var maxSize int64 = MaxFileSize << 20
	lmt := io.MultiReader(buffer, io.LimitReader(p, maxSize-511))
	written, err := io.Copy(f, lmt)
	if err != nil && err != io.EOF {
		return "", types.ErrBodyTooLarge
	}
	if written > maxSize {
		os.Remove(f.Name())
		return "", fmt.Errorf("Error while deleting file: %w", err)
	}
	return f.Name(), nil
}
