package stores

import (
	"api/models"
	"github.com/jmoiron/sqlx"
)

type userStoreInterface interface {
	Insert(u *models.User) error
	GetById(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByGithubId(id string) (*models.User, error)
	UpdateData(id, username, email, avatar string) error
	DeleteById(id string) error
}

type Users struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *Users {
	return &Users{
		db: db,
	}
}

func (s *Users) Insert(u *models.User) error {
	_, err := s.db.NamedExec(
		"insert into users(id, github_id, username, email, avatar, plan) values(:id, :github_id, :username, :email, :avatar, :plan)",
		u,
	)
	return err
}

func (s *Users) GetById(id string) (*models.User, error) {
	dbUser := new(models.User)
	err := s.db.Get(
		dbUser,
		"select * from users where id=$1",
		id,
	)
	return dbUser, err
}

func (s *Users) GetByEmail(email string) (*models.User, error) {
	dbUser := new(models.User)
	err := s.db.Get(
		dbUser,
		"select * from users where email=$1",
		email,
	)
	return dbUser, err
}

func (s *Users) GetByGithubId(id string) (*models.User, error) {
	dbUser := new(models.User)
	err := s.db.Get(
		dbUser,
		"select * from users where github_id=$1",
		id,
	)
	return dbUser, err
}

func (s *Users) UpdateData(id, username, email, avatar string) error {
	_, err := s.db.Exec(
		"update users set username=$1, email=$2, avatar=$3 where id=$4",
		username,
		email,
		avatar,
		id,
	)
	return err
}

func (s *Users) DeleteById(id string) error {
	_, err := s.db.Exec(
		"delete from users where id=$1",
		id,
	)
	return err
}
