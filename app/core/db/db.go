package db

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"helpers/app/bootstrap"
)

func CreatePoint(x, y float64) pgtype.Point {
	point := pgtype.Point{P: pgtype.Vec2{x, y}, Status: pgtype.Present}

	return point
}

type Conn interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

var pool *pgxpool.Pool

func Pool(ctx context.Context, cnf bootstrap.ConfigDb) *pgxpool.Pool {
	if pool == nil {
		pool = connect(ctx, cnf)
	}

	return pool
}

func GetPool() *pgxpool.Pool {
	return pool
}

func connect(ctx context.Context, cnf bootstrap.ConfigDb) *pgxpool.Pool {
	conf, err := pgxpool.ParseConfig(cnf.DSN())
	if err != nil {
		zap.S().Error(ctx, err)
	}

	p, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		zap.S().Errorf("Unable to connection to database: %v\n", err)
	}

	return p
}
