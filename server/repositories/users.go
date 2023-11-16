package repositories

import (
	"visio/models"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
  db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
  return &UserRepo{
    db: db,
  }
}

func (r UserRepo) InsertNewUser(user *models.User) error {
  _, err := r.db.NamedExec(
    "insert into users(id, github_id, username, avatar) values(:id, :github_id, :username, :avatar)",
    &user,
  )
  return err
}

func (r UserRepo) GetByGithubId(id string) (user *models.User, err error){
  user = &models.User{}
  err = r.db.Get(user, "select * from users where github_id=$1", id)
  return
}

func (r UserRepo) GetById(id string) (user *models.User, err error){
  user = new(models.User)
  err = r.db.Get(user, "select * from users where id=$1", id)
  return
}

func (r UserRepo) UpdateUserInfos(github_id , username , avatar string) error {
  _, err := r.db.Exec(
    "update users set username=$1, avatar=$2 where github_id=$3",
    username,
    avatar,
    github_id,
  )
  return err
}
