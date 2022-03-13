package soc_net

import "time"

type UserSocNet struct {
	ID      int       `json:"id" table_name:"users_soc_nets"`
	UserId  int64     `json:"user_id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Deleted time.Time `json:"deleted"`
}

func (s *UserSocNet) IsEntity() {}
