package core

import (
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func StrP(src string) *string {
	return &src
}

func IsRealError(err error) bool {
	res := err != nil && !errors.As(err, &pgx.ErrNoRows)

	return res
}
