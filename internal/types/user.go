package types

type User struct {
	Id           string `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	SignupDate   string `json:"signup_date" db:"signup_date"`
}
