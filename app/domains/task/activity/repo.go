package activity

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
	return core.NewTableSchema(&TaskActivity{})
}

func NewRepo(c db.Conn) *Repository {
	return &Repository{
		core.NewRepository(c, Schema),
	}
}

func (s *Repository) CreateOne(ctx context.Context, entity map[string]interface{}) (err error) {
	columns, vals := core.EntityToColumns(entity)
	err = core.CreateOne(ctx, s.Conn(), s.Schema().TableName(), columns, vals, nil)

	return
}

func (s *Repository) UpdateOne(ctx context.Context, entity map[string]interface{}, cond map[string]interface{}) (int, error) {
	return core.UpdateOne(ctx, s.ConnPool, s.Schema().TableName(), entity, cond)
}

func (s *Repository) FindOne(
	ctx context.Context,
	fields []string,
	cond map[string]interface{},
) (*TaskActivity, error) {
	res := &TaskActivity{}
	err := core.FindOne(ctx, s.ConnPool, s.Schema().TableName(), res, fields, cond)

	return res, err
}

func (s *Repository) FindMany(
	ctx context.Context,
	fields []string,
	cond map[string]interface{},
) ([]*TaskActivity, error) {
	var res []*TaskActivity
	err := core.FindMany(ctx, s.ConnPool, s.Schema().TableName(), &res, fields, cond)
	if err != nil {
		return nil, err
	}

	return res, nil
}
