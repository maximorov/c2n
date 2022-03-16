package executor

import (
	"github.com/jackc/pgx/pgtype"
)

type UserExecutor struct {
	ID       int          `json:"id" table_name:"users_executors"`
	UserId   int64        `json:"user_id"`
	Position pgtype.Point `json:"position"`
	Area     int          `json:"area"`
	City     string       `json:"city"`
	Inform   bool         `json:"inform"`
}

func (s *UserExecutor) IsEntity() {}
