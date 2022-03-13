package task

import (
	"context"
	"helpers/app/bootstrap"
	"helpers/app/core/db"
	"testing"
)

func TestService_CreateTask(t *testing.T) {
	bootstrap.InitEnv(`../../../`)
	bootstrap.InitConfig()
	bootstrap.InitLogger()

	ctx := context.Background()
	connPool := db.Pool(ctx, bootstrap.Cnf.DB)
	s := NewService(connPool)
	_, err := s.CreateTask(ctx, 1, 12, 13, `Светлана. Двое детей. Нужна еда.
Харьковское шоссе.
Номер: 099 222 33 44`)
	if err != nil {
		t.Error(err)
	}
}
