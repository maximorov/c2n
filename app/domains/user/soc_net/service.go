package soc_net

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

func (s *Service) CreateOne(ctx context.Context, userId int, userSocNetId string) (int, error) {
	id, err := s.Repo.CreateOne(ctx, map[string]interface{}{
		`user_id`:    userId,
		`soc_net_id`: userSocNetId,
	})

	return id, err
}

func (s *Service) GetOneByUserID(ctx context.Context, userID int) (*UserSocNet, error) {
	user, err := s.Repo.FindOne(ctx, []string{`id`, `soc_net_id`, `last_received_message`}, map[string]interface{}{
		`user_id`: userID,
	})

	return user, err
}

func (s *Service) GetOneBySocNetID(ctx context.Context, userSocNetID string) (*UserSocNet, error) {
	user, err := s.Repo.FindOne(ctx, []string{`id`, `user_id`}, map[string]interface{}{
		`soc_net_id`: userSocNetID,
	})

	return user, err
}
