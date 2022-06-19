// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const fooLimit = `-- name: FooLimit :many
SELECT a FROM foo
LIMIT $1
`

func (q *Queries) FooLimit(ctx context.Context, limit int32) ([]sql.NullString, error) {
	ctx, done := q.observer(ctx, "FooLimit")
	rows, err := q.db.QueryContext(ctx, fooLimit, limit)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var a sql.NullString
		if err := rows.Scan(&a); err != nil {
			return nil, done(err)
		}
		items = append(items, a)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}

const fooLimitOffset = `-- name: FooLimitOffset :many
SELECT a FROM foo
LIMIT $1 OFFSET $2
`

type FooLimitOffsetParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) FooLimitOffset(ctx context.Context, arg FooLimitOffsetParams) ([]sql.NullString, error) {
	ctx, done := q.observer(ctx, "FooLimitOffset")
	rows, err := q.db.QueryContext(ctx, fooLimitOffset, arg.Limit, arg.Offset)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var a sql.NullString
		if err := rows.Scan(&a); err != nil {
			return nil, done(err)
		}
		items = append(items, a)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}
