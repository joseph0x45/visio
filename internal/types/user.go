package types

import "time"

type User struct {
	Id         string    `json:"id" db:"id"`
	Email      string    `json:"email" db:"email"`
	Password   string    `json:"password" db:"password"`
	SignupDate time.Time `json:"signup_date" db:"signup_date"`
}
