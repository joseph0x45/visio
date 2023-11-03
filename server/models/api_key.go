package models

type Key struct {
	Id      string `json:"id" db:"id"`
	Owner   string `json:"owner" db:"owner"`
	Prefix  string `json:"prefix" db:"prefix"`
	KeyHash string `json:"key_hash" db:"key_hash"`
}
