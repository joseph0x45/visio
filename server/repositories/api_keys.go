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

func (r *Keys_repo) GetUserNumberOfKeys(user_id string) (count int, err error) {
  count = 0
  err = r.db.QueryRowx("select count(*) from keys where owner=$1", user_id).Scan(&count)
  return
}

func (r *Keys_repo) InsertNewKey(key *models.Key) error {
  _, err := r.db.NamedExec(
    "insert into keys(id, owner, prefix, key_hash) values(:id, :owner, :prefix, :key_hash)",
    &key,
  )
  return err
}

func (r *Keys_repo) DeleteKey(key_prefix string, user_id string) error {
  _, err := r.db.Exec("delete from keys where prefix=$1 and owner=$2", key_prefix, user_id)
  return err
}
