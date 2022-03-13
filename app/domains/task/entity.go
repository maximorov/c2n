package task

import (
	"helpers/app/core/db"
	"time"
)

type Task struct {
	ID        int `table_name:"tasks"`
	UserID    int
	Position  db.Point
	status    string
	text      string
	deadline  time.Time
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

func (s *Task) IsEntity() {}
