package executor

import (
	"helpers/app/core/db"
)

type UserExecutor struct {
	UserId   int64    `json:"user_id" table_name:"users_executors"`
	Position db.Point `json:"position"`
	Area     int      `json:"area"`
	City     string   `json:"city"`
}

func (s *UserExecutor) IsEntity() {}