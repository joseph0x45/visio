package types

import "errors"

var (
	ErrNoPrimaryEmailFound = errors.New("No primary email found")
	ErrUserNotFound         = errors.New("No user found")
	ErrSessionNotFound     = errors.New("Session not found")
)
