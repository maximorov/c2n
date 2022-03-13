package task

import (
	"context"
	"helpers/app/bootstrap"
	"helpers/app/core/db"
	"testing"
)

func TestRepository_UpdateOne(t *testing.T) {
	bootstrap.InitEnv(`../../../`)
	bootstrap.InitConfig()
	bootstrap.InitLogger()

	ctx := context.Background()
	connPool := db.Pool(ctx, bootstrap.Cnf.DB)
	s := NewRepo(connPool)
	_, err := s.UpdateOne(
		ctx,
		map[string]interface{}{`status`: `new`},
		map[string]interface{}{`id`: 1},
	)
	if err != nil {
		t.Error(err)
	}
}

func TestRepository_FindOne(t *testing.T) {
	bootstrap.InitEnv(`../../../`)
	bootstrap.InitConfig()
	bootstrap.InitLogger()

	ctx := context.Background()
	connPool := db.Pool(ctx, bootstrap.Cnf.DB)
	s := NewRepo(connPool)
	res, err := s.FindOne(
		ctx,
		[]string{`id`, `status`, `position`},
		map[string]interface{}{`id`: 1},
	)
	if err != nil || res == nil {
		t.Error(err)
	}
}

func TestRepository_FindMany(t *testing.T) {
	bootstrap.InitEnv(`../../../`)
	bootstrap.InitConfig()
	bootstrap.InitLogger()

	ctx := context.Background()
	connPool := db.Pool(ctx, bootstrap.Cnf.DB)
	s := NewRepo(connPool)
	res, err := s.FindMany(
		ctx,
		[]string{`id`, `status`, `position`},
		map[string]interface{}{`status`: `raw`},
	)
	if err != nil || res == nil {
		t.Error(err)
	}
}
