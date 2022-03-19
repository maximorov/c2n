package activity

import (
	"time"
)

var AllowedStatuses = map[string]map[string]bool{
	// TODO: write allowed statuses movement like in task.status
}

type TaskActivity struct {
	TaskID     int `json:"task_id" table_name:"tasks_activity"`
	ExecutorID int
	Status     string
	Deadline   time.Time
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty"`
}

func (s *TaskActivity) IsEntity() {}
