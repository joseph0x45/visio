package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

func NewPostgresPool() *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)
	return db
}
