// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"
)

const selectFoo = `-- name: SelectFoo :exec
SELECT a, b, c, d FROM foo
`

func (q *Queries) SelectFoo(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, selectFoo)
	return err
}
