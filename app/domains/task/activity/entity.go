package activity

import (
	"time"
)

type TaskActivity struct {
	ID         int `table_name:"tasks_activity"`
	ExecutorID int
	Status     string
	Deadline   time.Time
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty"`
}

func (s *TaskActivity) IsEntity() {}
