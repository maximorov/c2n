package usecase

import (
	"context"
	"helpers/app/core/db"
	"helpers/app/domains/task"
)

func NewTaskUseCase(connPool db.Conn) *TaskUseCase {
	return &TaskUseCase{taskRepo: task.NewRepo(connPool)}
}

type TaskUseCase struct {
	taskRepo *task.Repository
}

func (s *TaskUseCase) GetTasksForUser(ctx context.Context, userId int) ([]*task.Task, error) {
	res, err := s.taskRepo.FindMany(ctx, []string{
		`id`,
		`user_id`,
		`position`,
		`status`,
		`text`,
		`deadline`,
	}, map[string]interface{}{
		`status`: `new`,
	})

	return res, err
}
