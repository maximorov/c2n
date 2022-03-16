package task

import (
	"context"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"helpers/app/core/db"
	"math"
)

const PI float64 = 3.141592653589793

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

func (s *Service) GetUserUndoneTasks(ctx context.Context, userId int) ([]*Task, error) {
	tasks, err := s.repo.FindMany(ctx, []string{
		`id`, `status`, `text`,
	}, map[string]interface{}{
		`user_id`: userId,
		`status`:  []string{`new`, `in_progress`},
	})

	return tasks, err
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

func (s *Service) CountDistance(loc1, loc2 pgtype.Point) float64 {
	radlat1 := float64(PI * loc1.P.X / 180)
	radlat2 := float64(PI * loc2.P.X / 180)

	theta := float64(loc1.P.Y - loc2.P.Y)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515 * 1609.344

	return dist
}

func (s *Service) FindTasksInRadius(ctx context.Context, location pgtype.Point, area float64) ([]*Task, error) {
	var result []*Task
	tasks, err := s.repo.FindMany(ctx, []string{`id`, `position`}, map[string]interface{}{
		`status`: `new`,
	})
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		dist := s.CountDistance(task.Position, location)
		if dist < area {
			result = append(result, task)
		}

	}

	return result, nil
}
