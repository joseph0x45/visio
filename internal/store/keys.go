package store

import (
	"fmt"
	"visio/internal/types"

	"github.com/jmoiron/sqlx"
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
      insert into keys(key_owner, prefix, key_hash, key_creation_date)
      values (:key_owner, :prefix, :key_hash, :key_creation_date)
    `,
		key,
	)

	if err != nil {
		return fmt.Errorf("Error while inserting new key: %w", err)
	}

	return nil
}

func (k *Keys) CountByOwnerId(ownerId string) (int, error) {
	count := 0
	err := k.db.QueryRowx("select count(*) from keys where key_owner=$1", ownerId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("Error while counting keys by owner id: %w", err)
	}
	return count, nil
}
