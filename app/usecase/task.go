package usecase

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"helpers/app/core"
	"helpers/app/core/db"
	"helpers/app/domains/task"
	"helpers/app/domains/task/activity"
	"helpers/app/domains/user/soc_net"
)

const CantUpdateStatus = `Can't update task status`

func NewTaskUseCase(connPool db.Conn) *TaskUseCase {
	return &TaskUseCase{
		taskRepo:         task.NewRepo(connPool),
		taskActivityRepo: activity.NewRepo(connPool),
		socNetRepo:       soc_net.NewRepo(connPool),
	}
}

type TaskUseCase struct {
	taskRepo         *task.Repository
	taskActivityRepo *activity.Repository
	socNetRepo       *soc_net.Repository
}

func (s *TaskUseCase) GetTasksForUser(ctx context.Context, circle interface{}) ([]*task.Task, error) {
	res, err := s.taskRepo.FindMany(ctx, []string{
		`id`,
		`user_id`,
		`position`,
		`status`,
		`text`,
		`deadline`,
	}, map[string]interface{}{
		`status`: `new`,
	})

	return res, err
}

func (s *TaskUseCase) UpdateTaskStatus(ctx context.Context, taskId int, status string) error {
	currentTask, err := s.taskRepo.FindOne(
		ctx,
		[]string{`status`},
		map[string]interface{}{`id`: taskId})
	if err != nil {
		return err
	}

	// check if status can be changed
	var ok bool
	var allowedStatuses map[string]bool

	if allowedStatuses, ok = task.AllowedStatuses[currentTask.Status]; !ok {
		zap.S().Error(fmt.Sprintf(`No status '%s' in allowed status list`, currentTask.Status))
		return errors.New(CantUpdateStatus)
	}
	if _, ok = allowedStatuses[status]; !ok {
		return errors.New(CantUpdateStatus)
	}

	_, err = s.taskRepo.UpdateOne(ctx, map[string]interface{}{
		`status`: status,
	}, map[string]interface{}{
		`id`: taskId,
	})

	return err
}

func (s *TaskUseCase) CreateRawTask(ctx context.Context, userId int, x, y float64) error {
	p := db.CreatePoint(x, y)
	_, err := s.taskRepo.CreateOne(ctx, map[string]interface{}{
		`user_id`:  userId,
		`text`:     ``,
		`position`: &p,
	})

	return err
}

func (s *TaskUseCase) UpdateLastRawWithText(ctx context.Context, taskId int, text string) error {
	_, err := s.taskRepo.UpdateOne(ctx, map[string]interface{}{
		`text`:   text,
		`status`: `new`,
	},
		map[string]interface{}{
			`id`: taskId,
		})

	return err
}

func (s *TaskUseCase) GetSocUserByTask(ctx context.Context, taskId int) (*soc_net.UserSocNet, error) {
	taskUser, err := s.taskRepo.FindOne(
		ctx,
		[]string{`user_id`},
		map[string]interface{}{`id`: taskId})
	if err != nil {
		return nil, err
	}

	socNetUser, err := s.socNetRepo.FindOne(
		ctx,
		[]string{`soc_net_id`},
		map[string]interface{}{`user_id`: taskUser.UserID})
	if err != nil {
		return nil, err
	}

	return socNetUser, nil
}

func (s *TaskUseCase) GetSocExecutorByTaskActivity(ctx context.Context, taskId int) (*soc_net.UserSocNet, error) {
	taskExecutor, err := s.taskActivityRepo.FindOne(
		ctx,
		[]string{`executor_id`},
		map[string]interface{}{`task_id`: taskId, `status`: `taken`})
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	socNetUser, err := s.socNetRepo.FindOne(
		ctx,
		[]string{`soc_net_id`},
		map[string]interface{}{`user_id`: taskExecutor.ExecutorID})
	if err != nil {
		return nil, err
	}

	return socNetUser, nil
}

func (s *TaskUseCase) GetExecutorUndoneTasks(ctx context.Context, userId int) ([]*task.Task, error) {
	var tasks []*task.Task

	tName := s.taskRepo.Schema().TableName()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sb := psql.Select([]string{`id`, `t.status`, `text`}...).
		From(`"`+tName+`" as t`).
		InnerJoin(`"tasks_activity" as ta ON t.id = ta.task_id`).
		Where(`ta.executor_id = ?`, userId).
		Where(`ta.status = 'taken'`)

	err := core.FindManySB(ctx, s.taskRepo.Conn(), sb, &tasks)

	return tasks, err
}
