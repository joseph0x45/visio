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
    insert into users(id, email, password, signup_date)
    values (:id, :email, :password, :signup_date)
    `,
		user,
	)
	return fmt.Errorf("Error while inserting new user: %w", err)
}

func (s *Users) CountByEmail(email string) (int, error) {
	count := 0
	err := s.db.QueryRowx(
		"select count(*) from users where email=$1",
		email,
	).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("Error while counting users by email: %w", err)
	}
	return count, nil
}
