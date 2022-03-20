package task

import (
	"github.com/jackc/pgx/pgtype"
	"time"
)

const TaskDeadline = 24 // days

const StatusRaw = `raw`
const StatusNew = `new`
const StatusInProgress = `in_progress`
const StatusDone = `done`
const StatusExpired = `expired`
const StatusCancelled = `cancelled`
const StatusRefused = `refused`

var AllowedStatuses = map[string]map[string]bool{
	StatusRaw: {
		StatusNew:     true,
		StatusExpired: true,
	},
	StatusNew: {
		StatusInProgress: true,
		StatusCancelled:  true,
		StatusRefused:    true,
		StatusExpired:    true,
	},
	StatusInProgress: {
		StatusNew:       true,
		StatusDone:      true,
		StatusCancelled: true,
		StatusRefused:   true,
		StatusExpired:   true,
	},
	StatusDone: {},
	StatusExpired: {
		StatusNew:  true,
		StatusDone: true,
	},
	StatusCancelled: {
		StatusNew: true,
	},
	StatusRefused: {},
}

type Task struct {
	ID       int `table_name:"tasks"`
	UserID   int
	Position pgtype.Point
	Status   string
	Text     string
	Deadline time.Time
	Created  time.Time `json:"created,omitempty"`
	Updated  time.Time `json:"updated,omitempty"`
}

func (s *Task) IsEntity() {}
