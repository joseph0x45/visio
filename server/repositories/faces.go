package repositories

import "github.com/jmoiron/sqlx"

type FacesRepo struct {
  db *sqlx.DB
}

func NewFacesRepo(db *sqlx.DB) *FacesRepo {
  return &FacesRepo{
    db: db,
  }
}
