package store

import (
	"database/sql"
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
      insert into users(id, email, password_hash, signup_date)
      values (:id, :email, :password_hash, :signup_date)
    `,
		user,
	)
	if err != nil {
		return fmt.Errorf("Error while inserting new user: %w", err)
	}
	return nil
}

func (s *Users) GetById(id string) (*types.User, error) {
	dbUser := new(types.User)
	err := s.db.Get(dbUser, "select * from users where id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrUserNotFound
		}
		return nil, fmt.Errorf("Error while querying user from database by id: %w", err)
	}
	return dbUser, nil
}

func (s *Users) GetByEmail(email string) (*types.User, error) {
	dbUser := new(types.User)
	err := s.db.Get(dbUser, "select * from users where email=$1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrUserNotFound
		}
		return nil, fmt.Errorf("Error while querying user from database by email: %w", err)
	}
	return dbUser, nil
}

func (s *Users) CountByEmail(email string) (int, error) {
	count := 0
	err := s.db.QueryRowx("select count(*) from users where email=$1", email).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("Error while counting users by email: %w", err)
	}
	return count, nil
}
