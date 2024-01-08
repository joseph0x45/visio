package types

import "errors"

var (
	ErrUserNotFound    = errors.New("No user found")
	ErrSessionNotFound = errors.New("Session not found")
	ErrDuplicatePrefix = errors.New("Duplicate key prefix")
)
