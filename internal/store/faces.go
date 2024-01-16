package store

import (
	"database/sql"
	"fmt"
	"visio/internal/types"

	"github.com/jmoiron/sqlx"
)

type Faces struct {
	db *sqlx.DB
}

func NewFacesStore(db *sqlx.DB) *Faces {
	return &Faces{
		db: db,
	}
}

func (f *Faces) CountByLabel(label string) (int, error) {
	count := 0
	err := f.db.QueryRowx("select count(*) from faces where label=$1", label).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("Error while counting faces: %w", err)
	}
	return count, nil
}

func (f *Faces) Save(face *types.Face) error {
	_, err := f.db.NamedExec(
		`
      insert into faces(
        id, label, user_id, descriptor
      )
      values(
        :id, :label, :user_id, :descriptor
      )
    `,
		face,
	)
	if err != nil {
		return fmt.Errorf("Error while inserting face: %w", err)
	}
	return nil
}

func (f *Faces) GetById(id, userId string) (*types.Face, error) {
	face := new(types.Face)
	err := f.db.Get(
		face,
		"select * from faces where id=$1 and user_id=$2",
		id,
		userId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrFaceNotFound
		}
		return nil, fmt.Errorf("Error while querying face: %w", err)
	}
	return face, nil
}

func (f *Faces) GetByUserId(userId string) ([]*types.Face, error) {
	faces := []*types.Face{}
	err := f.db.Select(
		&faces,
		"select * from faces where user_id=$1",
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("Error while querying faces by user_id %w", err)
	}
	return faces, nil
}

func (f *Faces) Delete(id, userId string) error {
	_, err := f.db.Exec("delete from faces where id=$1 and user_id=$2", id, userId)
	if err != nil {
		return fmt.Errorf("Error while deleting face: %w", err)
	}
	return nil
}
