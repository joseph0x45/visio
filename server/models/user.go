package models

type User struct {
  Id string `json:"id" db:"id"`
  GithubId string `json:"github_id" db:"github_id"`
  Username string `json:"username" db:"username"`
  Email string `json:"email" db:"email"`
  Avatar string `json:"avatar" db:"avatar"`
  Plan string `json:"plan" db:"plan"`
}
