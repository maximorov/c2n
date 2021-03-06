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

var TranslatedStatuses = map[string]string{
	StatusRaw:        `не закінчено`,
	StatusNew:        `нове`,
	StatusInProgress: `виконується`,
	StatusDone:       `виконано`,
	StatusExpired:    `не актуально`,
	StatusCancelled:  `Видалено`,
	StatusRefused:    `заблоковано`,
}

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
	StatusDone: {
		StatusNew: true,
	},
	StatusExpired: {
		StatusNew:  true,
		StatusDone: true,
	},
	StatusCancelled: {
		StatusNew: true,
	},
	StatusRefused: {},
}

const TaskTextLength = 255

type Task struct {
	ID       int `table_name:"tasks"`
	UserID   int
	Position pgtype.Point
	Status   string
	Text     string
	Deadline time.Time
	Created  time.Time `json:"created,omitempty"`
	Updated  time.Time `json:"updated,omitempty"`

	distance float64
}

func (s *Task) SetDistance(dist float64) {
	s.distance = dist
}

func (s *Task) GetDistance() float64 {
	return s.distance
}

func (s *Task) IsEntity() {}
