package task

import (
	"context"
	"helpers/app/core"
)

func NewService(connPool core.Conn) *Service {
	return &Service{repo: NewRepo(connPool)}
}

type Service struct {
	repo *Repository
}

func (s *Service) CreateTask(ctx context.Context, text string) error {
	_, err := s.repo.CreateOne(ctx, map[string]interface{}{`text`: text})

	return err
}
