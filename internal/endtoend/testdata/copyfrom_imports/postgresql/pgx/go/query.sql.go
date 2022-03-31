// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const insertSingleValue = `-- name: InsertSingleValue :exec
INSERT INTO myschema.foo (a) VALUES ($1)
`

func (q *Queries) InsertSingleValue(ctx context.Context, a sql.NullString) error {
	_, err := q.db.Exec(ctx, insertSingleValue, a)
	return err
}

type InsertValuesParams struct {
	A sql.NullString
	B sql.NullInt32
}
