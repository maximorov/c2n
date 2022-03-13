package executor

import (
	"context"
	"helpers/app/core"
	"helpers/app/core/db"
)

var Schema = NewSchema()

type Repository struct {
	*core.Repository
}

func NewSchema() *core.TableSchema {
	return core.NewTableSchema(&UserExecutor{})
}

func NewRepo(c db.Conn) *Repository {
	return &Repository{
		core.NewRepository(c, Schema),
	}
}

func (s *Repository) CreateOne(ctx context.Context, entity map[string]interface{}) (retId int, err error) {
	columns, vals := core.EntityToColumns(entity)

	err = core.CreateOne(ctx, s.Conn(), s.Schema().TableName(), columns, vals, &retId)

	return
}

func (s *Repository) UpdateOne(ctx context.Context, entity map[string]interface{}, cond map[string]interface{}) (int, error) {
	return core.UpdateOne(ctx, s.ConnPool, s.Schema().TableName(), entity, cond)
}

func (s *Repository) FindOne(
	ctx context.Context,
	fields []string,
	cond map[string]interface{},
) (*UserExecutor, error) {
	res := &UserExecutor{}
	res.Position = db.CreatePoint(0, 0)
	err := core.FindOne(ctx, s.ConnPool, s.Schema().TableName(), res, fields, cond)

	return res, err
}

func (s *Repository) FindMany(
	ctx context.Context,
	fields []string,
	cond map[string]interface{},
) ([]*UserExecutor, error) {
	var res []*UserExecutor
	err := core.FindMany(ctx, s.ConnPool, s.Schema().TableName(), &res, fields, cond)
	if err != nil {
		return nil, err
	}

	return res, nil
}
