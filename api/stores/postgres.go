package stores

import (
	"os"
  _ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

func GetPostgresConnection() *sqlx.DB {
  db, err := sqlx.Connect("postgres", os.Getenv("DB_URL"))
  if err!= nil{
    panic(err)
  }
  return db
}
