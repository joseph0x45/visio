package models

type User struct {
	Id       string `db:"id"`
	GithubId string `db:"github_id"`
	Username string `db:"username"`
	Avatar   string `db:"avatar"`
	Plan     string `db:"plan"`
}

type UserData struct {
	Username string `json:"username" db:"username"`
	Avatar   string `json:"avatar" db:"avatar"`
	Plan     string `json:"plan" db:"plan"`
}
