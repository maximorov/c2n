package executor

import (
	"helpers/app/core/db"
)

type UserExecutor struct {
	ID       int      `json:"id" table_name:"users_executors"`
	UserId   int64    `json:"user_id"`
	Position db.Point `json:"position"`
	Area     int      `json:"area"`
	City     string   `json:"city"`
}

func (s *UserExecutor) IsEntity() {}
