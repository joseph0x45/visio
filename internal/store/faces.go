package store

import (
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
