package types

type Face struct {
	Id         string `json:"id" db:"id"`
	Label      string `json:"label" db:"label"`
	UserId     string `json:"user_id" db:"user_id"`
	Descriptor string `json:"descriptor" db:"descriptor"`
}
