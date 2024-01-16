package types

import "errors"

var (
	ErrFileNotFoundMessage      = "'face' field not found in request body"
	ErrUnsupportedFormatMessage = "File type not supported. Only image/jpeg and image/png are."
	ErrBodyTooLargeMessage      = "Request body too large. Only 5Mb is allowed."
)

var (
	ErrUserNotFound      = errors.New("No user found")
	ErrSessionNotFound   = errors.New("Session not found")
	ErrDuplicatePrefix   = errors.New("Duplicate key prefix")
	ErrFileNotFound      = errors.New(ErrFileNotFoundMessage)
	ErrUnsupportedFormat = errors.New(ErrUnsupportedFormatMessage)
	ErrBodyTooLarge      = errors.New(ErrBodyTooLargeMessage)
	ErrKeyNotFound       = errors.New("Key not found")
	ErrFaceNotFound      = errors.New("Face not found")
)
