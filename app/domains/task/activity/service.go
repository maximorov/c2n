package activity

import (
	"context"
	"helpers/app/core/db"
)

func NewService(connPool db.Conn) *Service {
	return &Service{repo: NewRepo(connPool)}
}

type Service struct {
	repo *Repository
}

func (s *Service) CreateActivity(ctx context.Context, userId, taskId int, status string) error {
	err := s.repo.CreateOne(ctx, map[string]interface{}{
		`executor_id`: userId,
		`task_id`:     taskId,
		`status`:      status,
	})

	return err
}
