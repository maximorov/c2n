package activity

import (
	"context"
	"helpers/app/core/db"
)

func NewService(connPool db.Conn) *Service {
	return &Service{Repo: NewRepo(connPool)}
}

type Service struct {
	Repo *Repository
}

func (s *Service) CreateActivity(ctx context.Context, userId, taskId int, status string) error {
	err := s.Repo.CreateOne(ctx, map[string]interface{}{
		`executor_id`: userId,
		`task_id`:     taskId,
		`status`:      status,
	})

	return err
}

func (s *Service) UpdateActivity(ctx context.Context, userId, taskId int, status string) error {
	_, err := s.Repo.UpdateOne(ctx, map[string]interface{}{
		`status`: status,
	}, map[string]interface{}{
		`executor_id`: userId,
		`task_id`:     taskId,
	})

	return err
}
