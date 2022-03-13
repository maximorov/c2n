package task

import (
	"helpers/app/core"
	"time"
)

type Task struct {
	ID        int `table_name:"tasks"`
	UserID    int
	Location  core.Point
	status    string
	text      string
	deadline  time.Time
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

func (s *Task) IsEntity() {}
