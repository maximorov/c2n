package executor

type UserExecutor struct {
	UserId int64       `json:"user_id" table_name:"users_executors"`
	Area   interface{} `json:"area"`
	City   string      `json:"city"`
}

func (s *UserExecutor) IsEntity() {}
