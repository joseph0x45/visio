package types

type Key struct {
	Id           string `json:"id" db:"id"`
	UserId       string `json:"user_id" db:"user_id"`
	Prefix       string `json:"prefix" db:"prefix"`
	KeyHash      string `json:"key_hash" db:"key_hash"`
	CreationDate string `json:"creation_date" db:"creation_date"`
}
