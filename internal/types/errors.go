package types

import "errors"

var (
	ErrUserNotFound    = errors.New("No user found")
	ErrSessionNotFound = errors.New("Session not found")
	ErrDuplicatePrefix = errors.New("Duplicate key prefix")
	ErrFileNotFound      = errors.New("File field not found in request body")
	ErrUnsupportedFormat = errors.New("Field type not supported")
)
