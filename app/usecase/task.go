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

func (s *TaskUseCase) GetTasksForUser(ctx context.Context, circle interface{}) ([]*task.Task, error) {
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

func (s *TaskUseCase) CreateRawTask(ctx context.Context, userId int, x, y float64) error {
	_, err := s.taskRepo.CreateOne(ctx, map[string]interface{}{
		`user_id`:  userId,
		`text`:     ``,
		`position`: db.CreatePoint(x, y),
	})

	return err
}

func (s *TaskUseCase) UpdateLastRawWithText(ctx context.Context, taskId int, text string) error {
	_, err := s.taskRepo.UpdateOne(ctx, map[string]interface{}{
		`text`:   text,
		`status`: `new`,
	},
		map[string]interface{}{
			`id`: taskId,
		})

	return err
}
