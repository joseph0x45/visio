package store

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"visio/internal/types"
)

type Users struct {
	db *sqlx.DB
}

func NewUsersStore(db *sqlx.DB) *Users {
	return &Users{
		db: db,
	}
}

func (s *Users) Insert(user *types.User) error {
	_, err := s.db.NamedExec(
		`
    insert into users(id, github_id, email, username, avatar, credits, signup_date)
    values (:id, :github_id, :email, :username, :avatar, :credits, :signup_date)
    `,
		user,
	)
	return fmt.Errorf("Error while inserting new user: %w", err)
}
