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

func (r *FacesRepo) InsertFace(face *models.Face) error {
  _, err := r.db.NamedExec(
    "insert into faces(id, created_by, descriptor, created_at, last_updated) values(:id, :created_by, :descriptor, :created_at, :last_updated)",
    &face,
  )
  return err
}

func (r *FacesRepo) SelectAllFacesCreatedByUser(user_id string) (faces []models.Face, err error) {
	faces = []models.Face{}
	err = r.db.Select(&faces, "select * from faces where created_by=$1", user_id)
	return
}

func (r *FacesRepo) DeleteFace(face_id , user_id string) error {
  _, err := r.db.Exec("delete from faces where id=$1 and created_by=$2", face_id, user_id)
  return err
}

func (r *FacesRepo) GetFaceById(face_id , user_id string) (face *models.Face, err error) {
  err = r.db.Get(
    face,
    "select * from faces where id=$1 and created_by=$2",
    face_id,
    user_id,
  )
  return
}

func (r *FacesRepo) UpdateFace(face_id , user_id, descriptor, last_updated string) error {
  _, err := r.db.Exec("update faces set descriptor=$1, last_updated=$2 where id=$3 and created_by=$4", descriptor, last_updated, face_id, user_id)
  return err
}
