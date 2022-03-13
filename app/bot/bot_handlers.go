package bot

import (
	"context"
	"go.uber.org/zap"
	"helpers/app/domains/task/activity"
	"strconv"
	"strings"
)

type Handler struct {
	TaskActivityService *activity.Service
}

func (s *Handler) Handle(msgText string) bool {
	handled := false

	switch {
	case strings.Contains(msgText, `hide:`) || strings.Contains(msgText, `accept:`):
		handled = true
		parsed := strings.Split(msgText, `:`)
		action := parsed[0]
		taskId, _ := strconv.Atoi(parsed[1])

		switch action {
		case `accept`:
			err := s.TaskActivityService.CreateActivity(context.TODO(), 1, taskId, `taken`)
			if err != nil {
				zap.S().Error(err)
			}
		case `hide`:
			err := s.TaskActivityService.CreateActivity(context.TODO(), 1, taskId, `hidden`)
			if err != nil {
				zap.S().Error(err)
			}
		}
	}

	return handled
}
