package core

import "github.com/jackc/pgx/pgtype"

type Point interface {
	pgtype.Value
	pgtype.BinaryDecoder
	pgtype.TextDecoder
	pgtype.BinaryEncoder
}

func CreatePoint(x, y float64) Point {
	point := &pgtype.Point{P: pgtype.Vec2{x, y}, Status: pgtype.Present}

	return point
}
