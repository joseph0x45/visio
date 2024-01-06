package database

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func connect() (db *sqlx.DB, err error) {
	db, err = sqlx.Connect("postgres", os.Getenv("PG_URL"))
	return
}

func NewPostgresPool() *sqlx.DB {
	var db *sqlx.DB
	db, err := connect()
	if err != nil {
		fmt.Println("Failed to connect to database server. Retrying in 5 seconds")
		time.Sleep(time.Second * 5)
		db, err = connect()
		if err != nil {
			fmt.Println("Failed to connect to database")
			panic(err)
		}
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)
	return db
}
