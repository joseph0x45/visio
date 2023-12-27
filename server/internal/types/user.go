package types

import "time"

type User struct {
	Id         string    `json:"id" db:"id"`
	GithubId   string    `json:"github_id" db:"github_id"`
	Username   string    `json:"username" db:"username"`
	Avatar     string    `json:"avatar" db:"avatar"`
	SignupDate time.Time `json:"signup_date" db:"signup_date"`
}
