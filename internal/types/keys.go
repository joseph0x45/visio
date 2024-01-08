package types

import "time"

type Key struct {
	UserId       string    `json:"user_id" db:"user_id"`
	Prefix       string    `json:"prefix" db:"prefix"`
	KeyHash      string    `json:"key_hash" db:"key_hash"`
	CreationDate time.Time `json:"creation_date" db:"creation_date"`
}
