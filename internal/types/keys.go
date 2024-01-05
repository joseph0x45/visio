package types

import "time"

type Key struct {
	KeyOwner           string    `json:"key_owner" db:"key_owner"`
	Prefix          string    `json:"prefix" db:"prefix"`
	KeyHash             string    `json:"key_hash" db:"key_hash"`
	KeyCreationDate time.Time `json:"key_creation_date" db:"key_creation_date"`
}
