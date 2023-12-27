package types

import "errors"

var (
	ErrNoPrimaryEmailFound = errors.New("No primary email found")
	ErrNoUserFound         = errors.New("No user found")
	ErrSessionNotFound     = errors.New("Session not found")
)
