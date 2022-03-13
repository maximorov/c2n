package soc_net

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
	return core.NewTableSchema(&UserSocNet{})
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
