package user

import "time"

type Role string

const Unknown Role = `unknown`
const Executor Role = `executor`
const Needy Role = `needy`

type User struct {
	ID          int       `json:"id" table_name:"users"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Deleted     time.Time `json:"deleted"`

	role Role
}

func (s *User) SetRole(role Role) {
	s.role = role
}

func (s *User) GetRole() Role {
	return s.role
}

func (s *User) IsEntity() {}
