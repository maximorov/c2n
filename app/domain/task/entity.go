package task

import (
	"time"
)

type Task struct {
	ID uint
	//constraint tasks_pk
	//primary key,
	UserID    uint
	Location  float64
	status    string
	text      string
	deadline  time.Time
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

func (s *Task) IsEntity() {}
