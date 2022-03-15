package soc_net

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

func (s *Service) CreateOne(ctx context.Context, userId int, userSocNetId int) (int, error) {
	id, err := s.repo.CreateOne(ctx, map[string]interface{}{
		`user_id`:    userId,
		`soc_net_id`: userSocNetId,
	})

	return id, err
}

func (s *Service) GetOneBySocNetID(ctx context.Context, userSocNetID int) (*UserSocNet, error) {
	user, err := s.repo.FindOne(ctx, []string{`id`, `user_id`}, map[string]interface{}{
		`soc_net_id`: userSocNetID,
	})

	return user, err
}
