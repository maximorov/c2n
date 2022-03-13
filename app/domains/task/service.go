package task

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

func (s *Service) CreateTask(ctx context.Context, userId int, x, y float64, text string) (int, error) {
	id, err := s.repo.CreateOne(ctx, map[string]interface{}{
		`user_id`:  userId,
		`text`:     text,
		`position`: db.CreatePoint(x, y),
	})

	return id, err
}
