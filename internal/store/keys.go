package store

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"visio/internal/types"
)

type Keys struct {
	db *sqlx.DB
}

func NewKeysStore(db *sqlx.DB) *Keys {
	return &Keys{
		db: db,
	}
}

func (k *Keys) Insert(key *types.Key, tx *sqlx.Tx) error {
	_, err := tx.NamedExec(
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

func (k *Keys) Delete(tx *sqlx.Tx, userId string) error {
	_, err := tx.Exec("delete from keys where user_id=$1", userId)
	if err != nil {
		return fmt.Errorf("Error while deleting key: %w", err)
	}
	return nil
}
