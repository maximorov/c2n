package task

import (
	"time"
)

const TaskDeadline = 24 // days

type Task struct {
	ID        int `table_name:"tasks"`
	UserID    int
	Position  interface{}
	Status    string
	Text      string
	Deadline  time.Time
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

func (s *Task) IsEntity() {}
