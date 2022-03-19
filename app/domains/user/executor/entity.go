package executor

import (
	"github.com/jackc/pgx/pgtype"
)

type UserExecutor struct {
	UserId   int          `json:"user_id" table_name:"users_executors"`
	Position pgtype.Point `json:"position"`
	Area     int          `json:"area"`
	City     string       `json:"city"`
	Inform   bool         `json:"inform"`
}

func (s *UserExecutor) IsEntity() {}
