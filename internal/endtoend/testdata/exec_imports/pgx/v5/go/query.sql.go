// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: query.sql

package querytest

import (
	"context"
)

const bar = `-- name: Bar :exec
SELECT bar
FROM foo
`

func (q *Queries) Bar(ctx context.Context) error {
	_, err := q.db.Exec(ctx, bar)
	return err
}

const bars = `-- name: Bars :exec
SELECT bars
FROM foo
`

func (q *Queries) Bars(ctx context.Context) error {
	_, err := q.db.Exec(ctx, bars)
	return err
}
