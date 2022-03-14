package user

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

func (s *Service) CreateOne(ctx context.Context, name string, phoneNumber int) (int, error) {
	id, err := s.repo.CreateOne(ctx, map[string]interface{}{
		`name`:         name,
		`phone_number`: phoneNumber,
	})

	return id, err
}

func (s *Service) UpdateOne(ctx context.Context, entity map[string]interface{}, cond map[string]interface{}) (int, error) {
	return s.repo.UpdateOne(ctx, entity, cond)
}

func (s *Service) GetOneByID(ctx context.Context, userId int) (*User, error) {
	user, err := s.repo.FindOne(ctx, []string{`id`, `phone_number`}, map[string]interface{}{
		`id`: userId,
	})

	return user, err
}
