package models

import "time"

type Face struct {
	Id          string    `json:"id" db:"id"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	Descriptor  string    `json:"descriptor" db:"descriptor"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	LastUpdated time.Time `json:"last_updated" db:"last_updated"`
}
