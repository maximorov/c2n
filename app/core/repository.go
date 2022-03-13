package core

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"helpers/app/core/db"
	"reflect"
	"time"
)

type Entity interface {
	IsEntity()
}

type deleteType int

const (
	softDelete deleteType = 0
	hardDelete deleteType = 1
)

type TableSchema struct {
	entity     Entity
	deleteType deleteType

	_table string
}

func NewTableSchema(e Entity) *TableSchema {
	res := &TableSchema{
		entity:     e,
		deleteType: hardDelete,
	}

	res.cacheTableName()

	return res
}

func (s *TableSchema) cacheTableName() {
	v := reflect.TypeOf(s.entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	f := v.Field(0)
	val, ok := f.Tag.Lookup("table_name")
	if !ok {
		zap.S().Fatalf("Table name for entity %#v not defined", s.entity)
	}

	s._table = val
}

func (s *TableSchema) TableName() string { // @TODO: Add cache
	return s._table
}

func NewRepository(cs db.Conn, sch *TableSchema) *Repository {
	return &Repository{cs, sch}
}

type Repository struct {
	ConnPool db.Conn
	schema   *TableSchema
}

func (s *Repository) SetSchema(sch *TableSchema) {
	s.schema = sch
}

func (s *Repository) Schema() *TableSchema {
	return s.schema
}

func (s *Repository) Conn() db.Conn {
	return s.ConnPool
}

func CreateOne(ctx context.Context, conn db.Conn, tName string, columns []string, vals []interface{}, retId interface{}) (err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Insert(`"` + tName + `"`).
		Columns(columns...).
		Values(vals...)

	if retId != nil {
		sql, _, _ := query.ToSql()
		err = conn.QueryRow(ctx, sql+" RETURNING id", vals...).Scan(retId)
	} else {
		sql, _, _ := query.ToSql()
		_, err = conn.Exec(ctx, sql, vals...)
	}

	return
}

func UpdateOne(ctx context.Context, conn db.Conn, tName string, entity map[string]interface{}, cond map[string]interface{}) (int, error) {
	lenEntity := len(entity)
	if lenEntity == 0 {
		return 0, nil
	}

	copyEntity := make(map[string]interface{}, lenEntity)
	for key, val := range entity {
		copyEntity[key] = val
	}

	if len(copyEntity) == 0 {
		return 0, nil
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	updateQuery := psql.Update(tName).Where(cond)

	for col, val := range copyEntity {
		updateQuery = updateQuery.Set(col, val)
	}

	sql, args, err := updateQuery.ToSql()
	if err != nil {
		return 0, err
	}

	cmdTag, err := conn.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}

	return int(cmdTag.RowsAffected()), nil
}

// Gets one record by condition
func FindOne(ctx context.Context, conn db.Conn, tName string, result Entity, fields []string, condition map[string]interface{}) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sb := psql.Select(fields...).
		From(`"` + tName + `"`).
		Where(condition).
		Limit(1)

	sql, args, err := sb.ToSql()

	err = pgxscan.Get(ctx, conn, result, sql, args...)
	if err != nil {
		return errors.WithMessagef(err, "sql query %s, args %+v", sql, args)
	}
	return nil
}

// Gets all record by condition
// FindMany Gets many records by condition
func FindMany(ctx context.Context, conn db.Conn, tName string, result interface{}, fields []string, condition map[string]interface{}) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sb := psql.Select(fields...).
		From(`"` + tName + `"`).
		Where(condition)

	sql, args, err := sb.ToSql()

	err = pgxscan.Select(ctx, conn, result, sql, args...)
	if err != nil {
		return errors.WithMessagef(err, "sql query %s, args %+v", sql, args)
	}
	return nil
}

// entityToColumns reflects on a struct and returns the values of fields with `db` tags,
// or a map[string]interface{} and returns the keys.
func EntityToColumns(values interface{}) ([]string, []interface{}) {
	var field string

	v := reflect.ValueOf(values)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	fields := []string{}
	vals := []interface{}{}
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field = v.Type().Field(i).Tag.Get("json")
			if field != "" /*&& IndexOf(field, expect) != -1 */ {
				fields = append(fields, field)
				vals = append(vals, v.Field(i).Interface())
			}
		}
		return fields, vals
	} /*&& IndexOf(field, expect) != -1 */
	if v.Kind() == reflect.Map {
		for _, keyv := range v.MapKeys() {
			if keyv.String() != "" {
				fields = append(fields, keyv.String())
				vals = append(vals, v.MapIndex(keyv).Interface())
			}
		}
		return fields, vals
	}
	panic(fmt.Errorf("EntityToColumns requires a struct or a map, found: %s", v.Kind().String()))
}

// AddCurrentTimeIfNotSet add time only when field is not set
func AddCurrentTimeIfNotSet(ctx context.Context, fields []string, values []interface{}, timeFields ...string) ([]string, []interface{}) {
	now := time.Now()

	for _, field := range timeFields {
		if i := IndexOfDirectSearch(field, fields); i == -1 {
			fields = append(fields, field)
			values = append(values, now)
		} else {
			if tmp, ok := values[i].(time.Time); ok && tmp.IsZero() {
				values[i] = now
			} else if tmp, ok := values[i].(*time.Time); ok && (tmp == nil || tmp.IsZero()) {
				values[i] = &now
			}
		}
	}

	return fields, values
}

// IndexOfDirectSearch searches for needle in haystack slice of strings
// and returns the index or -1 if needle is not present in haystack.
func IndexOfDirectSearch(needle string, haystack []string) int {
	for i, s := range haystack {
		if needle == s {
			return i
		}
	}

	return -1
}
