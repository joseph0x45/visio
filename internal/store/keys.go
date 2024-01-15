package store

import (
	"database/sql"
	"fmt"
	"visio/internal/types"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Keys struct {
	db *sqlx.DB
}

func NewKeysStore(db *sqlx.DB) *Keys {
	return &Keys{
		db: db,
	}
}

func (k *Keys) Insert(key *types.Key) error {
	_, err := k.db.NamedExec(
		`
      insert into keys(id, user_id, prefix, key_hash, creation_date)
      values (:id, :user_id, :prefix, :key_hash, :creation_date)
    `,
		key,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				return types.ErrDuplicatePrefix
			}
		}
		return fmt.Errorf("Error while inserting new key: %w", err)
	}
	return nil
}

func (k *Keys) CountByOwnerId(userId string) (int, error) {
	count := 0
	err := k.db.QueryRowx("select count(*) from keys where user_id=$1", userId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("Error while counting keys by owner id: %w", err)
	}
	return count, nil
}

func (k *Keys) GetByPrefix(prefix string) (*types.Key, error) {
	key := new(types.Key)
	err := k.db.Get(
		key,
		"select * from keys where prefix=$1",
		prefix,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrKeyNotFound
		}
		return nil, fmt.Errorf("Error while retrieving key by prefix: %w", err)
	}
	return key, nil
}

func (k *Keys) GetByUserId(id string) ([]types.Key, error) {
	data := []types.Key{}
	err := k.db.Select(
		&data,
		`select * from keys where user_id=$1`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving keys from database: %w", err)
	}
	return data, nil
}

func (k *Keys) Delete(prefix, userId string) error {
	_, err := k.db.Exec("delete from keys where prefix=$1 and user_id=$2", prefix, userId)
	if err != nil {
		return fmt.Errorf("Error while deleting key: %w", err)
	}
	return nil
}
