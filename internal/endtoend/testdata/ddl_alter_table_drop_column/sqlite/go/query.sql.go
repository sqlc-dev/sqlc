// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
)

const placeholder = `-- name: Placeholder :exec
SELECT baz from foo
`

func (q *Queries) Placeholder(ctx context.Context) error {
	ctx, done := q.observer(ctx, "Placeholder")
	_, err := q.db.ExecContext(ctx, placeholder)
	return done(err)
}
