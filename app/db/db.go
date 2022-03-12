package db

import (
	"context"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"go.uber.org/zap"
	"helpers/app/bootstrap"
)

type (
	Txer interface {
		// Begin starts a pseudo nested transaction.
		Begin(ctx context.Context) (pgx.Tx, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
		Conn
	}
	Conn interface {
		Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
		Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
		QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	}
)

var pool *pgx.ConnPool

func Pool(ctx context.Context, cnf bootstrap.ConfigDb) *pgx.ConnPool {
	if pool == nil {
		pool = connect(ctx, cnf)
	}

	return pool
}

func connect(ctx context.Context, cnf bootstrap.ConfigDb) *pgx.ConnPool {
	conf, err := pgx.ParseDSN(cnf.DSN())
	if err != nil {
		zap.S().Error(ctx, err)
	}

	p, err := pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: conf})
	if err != nil {
		zap.S().Errorf("Unable to connection to database: %v\n", err)
	}

	return p
}
