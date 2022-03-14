package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"helpers/app/bootstrap"
)

type Point interface {
	pgtype.Value
	pgtype.BinaryDecoder
	pgtype.TextDecoder
	pgtype.BinaryEncoder
	pgtype.TextEncoder
	sql.Scanner
	driver.Valuer
}

func CreatePoint(x, y float64) Point {
	point := &pgtype.Point{P: pgtype.Vec2{x, y}, Status: pgtype.Present}

	return point
}

type Circle interface {
	pgtype.Value
	pgtype.BinaryDecoder
	pgtype.TextDecoder
	pgtype.BinaryEncoder
	pgtype.TextEncoder
	sql.Scanner
	driver.Valuer
}

func CreateCircle(x, y, r float64) Circle {
	point := &pgtype.Circle{P: pgtype.Vec2{x, y}, R: r, Status: pgtype.Present}

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
