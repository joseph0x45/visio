package types

import "time"

type Key struct {
	Owner           string    `json:"owner" db:"owner"`
	Prefix          string    `json:"prefix" db:"prefix"`
	Key             string    `json:"key" db:"key"`
	KeyCreationDate time.Time `json:"keyCreationDate" db:"key_creation_date"`
}
