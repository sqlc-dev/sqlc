// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const fooLimit = `-- name: FooLimit :many
SELECT a FROM foo
LIMIT ?
`

func (q *Queries) FooLimit(ctx context.Context, limit int64) ([]sql.NullString, error) {
	rows, err := q.db.QueryContext(ctx, fooLimit, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var a sql.NullString
		if err := rows.Scan(&a); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fooLimitOffset = `-- name: FooLimitOffset :many
SELECT a FROM foo
LIMIT ? OFFSET ?
`

type FooLimitOffsetParams struct {
	Limit  int64
	Offset int64
}

func (q *Queries) FooLimitOffset(ctx context.Context, arg FooLimitOffsetParams) ([]sql.NullString, error) {
	rows, err := q.db.QueryContext(ctx, fooLimitOffset, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var a sql.NullString
		if err := rows.Scan(&a); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
