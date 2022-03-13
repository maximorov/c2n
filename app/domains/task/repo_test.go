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
