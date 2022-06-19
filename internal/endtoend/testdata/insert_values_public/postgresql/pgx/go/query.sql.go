// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const insertValues = `-- name: InsertValues :exec
INSERT INTO public.foo (a, b) VALUES ($1, $2)
`

type InsertValuesParams struct {
	A sql.NullString
	B sql.NullInt32
}

func (q *Queries) InsertValues(ctx context.Context, arg InsertValuesParams) error {
	ctx, done := q.observer(ctx, "InsertValues")
	_, err := q.db.Exec(ctx, insertValues, arg.A, arg.B)
	return done(err)
}
