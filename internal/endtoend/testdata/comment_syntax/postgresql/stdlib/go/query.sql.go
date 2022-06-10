// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const doubleDash = `-- name: DoubleDash :one
SELECT bar FROM foo LIMIT 1
`

func (q *Queries) DoubleDash(ctx context.Context) (sql.NullString, error) {
	row := q.db.QueryRowContext(ctx, doubleDash)
	var bar sql.NullString
	err := row.Scan(&bar)
	return bar, err
}

const slashStar = `-- name: SlashStar :one
SELECT bar FROM foo LIMIT 1
`

func (q *Queries) SlashStar(ctx context.Context) (sql.NullString, error) {
	row := q.db.QueryRowContext(ctx, slashStar)
	var bar sql.NullString
	err := row.Scan(&bar)
	return bar, err
}
