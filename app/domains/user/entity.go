package user

import "time"

type User struct {
	ID          int       `json:"id" table_name:"tasks"`
	PhoneNumber int64     `json:"phone_number"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Deleted     time.Time `json:"deleted"`
}

func (s *User) IsEntity() {}
