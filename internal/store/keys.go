package store

import "github.com/jmoiron/sqlx"

type Keys struct {
	db *sqlx.DB
}

func NewKeysStore(db *sqlx.DB) *Keys {
	return &Keys{
		db: db,
	}
}
