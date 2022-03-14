package task

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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

func (s *Service) GetUsersRawTask(ctx context.Context, userId int) (*Task, error) {
	task, err := s.repo.FindOne(ctx, []string{
		`id`,
	}, map[string]interface{}{
		`user_id`: userId,
		`status`:  `raw`,
	})

	return task, err
}

func (s *Service) IsUserHaveUndoneTasks(ctx context.Context, userId int) bool {
	_, err := s.repo.FindOne(ctx, []string{
		`id`,
	}, map[string]interface{}{
		`user_id`: userId,
		`status`:  []string{`new`, `in_progress`},
	})
	if err != nil {
		if !errors.As(err, &pgx.ErrNoRows) {
			zap.S().Error(err)
		}
		return false
	}

	return true
}
