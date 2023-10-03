// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const insertValues = `-- name: InsertValues :exec
INSERT INTO public.foo (a, b) VALUES (?, ?)
`

type InsertValuesParams struct {
	A sql.NullString
	B sql.NullInt32
}

func (q *Queries) InsertValues(ctx context.Context, arg InsertValuesParams) error {
	query := insertValues
	queryParams := []interface{}{arg.A, arg.B}

	_, err := q.db.ExecContext(ctx, query, queryParams...)
	return err
}
