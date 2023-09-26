// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: query.sql

package db

import (
	"context"
)

const doStuff = `-- name: DoStuff :one
DO $$
    BEGIN
        ALTER TABLE authors
        ADD COLUMN marked_for_processing bool;
    END
$$
`

type DoStuffRow struct {
}

func (q *Queries) DoStuff(ctx context.Context) (DoStuffRow, error) {
	row := q.db.QueryRowContext(ctx, doStuff)
	var i DoStuffRow
	err := row.Scan()
	return i, err
}
