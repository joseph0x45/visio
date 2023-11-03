package repositories

import "github.com/jmoiron/sqlx"

type UserRepo struct {
  db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
  return &UserRepo{
    db: db,
  }
}

func (r UserRepo) Insert(username string){
  r.db.Exec("somethihg")
}
