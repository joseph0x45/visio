package store

import (
	"database/sql"
	"fmt"
	"visio/internal/types"

	"github.com/jmoiron/sqlx"
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
      insert into users(id, github_id, email, username, avatar, signup_date)
      values (:id, :github_id, :email, :username, :avatar, :signup_date)
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

func (s *Users) GetByGithubId(id string) (*types.User, error) {
	dbUser := new(types.User)
	err := s.db.Get(dbUser, "select * from users where github_id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrUserNotFound
		}
		return nil, fmt.Errorf("Error while querying user from database: %w", err)
	}
	return dbUser, nil
}

func (s *Users) UpdateUserData(id, email, username, avatar string) error {
	_, err := s.db.Exec(
		"update users set username=$1, email=$2, avatar=$3 where id=$4",
		username, email, avatar, id,
	)
	if err != nil {
		return fmt.Errorf("Error while updating user data: %w", err)
	}
	return nil
}
