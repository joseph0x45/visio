package repositories

import (
	"visio/models"

	"github.com/jmoiron/sqlx"
)

type Keys_repo struct {
  db *sqlx.DB
}

func NewKeysRepo(db *sqlx.DB) *Keys_repo {
  return &Keys_repo{
    db: db,
  }
}

func (r *Keys_repo) SelectKeys(user_id string) (keys []models.Key, err error) {
  keys = []models.Key{}
  err = r.db.Select(&keys, "select * from keys where owner=$1", user_id)
  return
}

func (r *Keys_repo) InsertNewKey(key *models.Key) error {
  _, err := r.db.NamedExec(
    "insert into keys(id, owner, prefix, key_hash) values(:id, :owner, :prefix, :key_hash)",
    &key,
  )
  return err
}
