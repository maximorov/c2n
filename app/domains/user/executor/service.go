package executor

import (
	"context"
	"github.com/jackc/pgx/pgtype"
	"helpers/app/core/db"
)

func NewService(connPool db.Conn) *Service {
	return &Service{repo: NewRepo(connPool)}
}

type Service struct {
	repo *Repository
}

func (s *Service) CreateOne(ctx context.Context, userId, area int, city string, pos pgtype.Point) (int, error) {
	id, err := s.repo.CreateOne(ctx, map[string]interface{}{
		`user_id`:  userId,
		`area`:     area,
		`city`:     city,
		`position`: pos,
	})

	return id, err
}

func (s *Service) UpdateOne(ctx context.Context, entity map[string]interface{}, cond map[string]interface{}) (int, error) {
	return s.repo.UpdateOne(ctx, entity, cond)
}

func (s *Service) GetOneByUserID(ctx context.Context, userId int) (*UserExecutor, error) {
	user, err := s.repo.FindOne(ctx, []string{`id`, `area`, `city`}, map[string]interface{}{
		`user_id`: userId,
	})

	return user, err
}
