package repositories

import (
	"visio/models"

	"github.com/jmoiron/sqlx"
)

type FacesRepo struct {
	db *sqlx.DB
}

func NewFacesRepo(db *sqlx.DB) *FacesRepo {
	return &FacesRepo{
		db: db,
	}
}

func (r *FacesRepo) SelectAllFacesCreatedByUser(user_id string) (faces []models.Face, err error) {
	faces = []models.Face{}
	err = r.db.Select(&faces, "select * from faces where created_by=$1", user_id)
	return
}

func (r *FacesRepo) DeleteFace(user_id , face_id string) (affected int64, err error) {
  res, err := r.db.Exec("delete from faces where id=$1 and created_by=$2", face_id, user_id)
  affected, _ = res.RowsAffected()
  return
}
