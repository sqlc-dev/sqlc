// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const barExists = `-- name: BarExists :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            bar
        where
            id = ?
    )
`

func (q *Queries) BarExists(ctx context.Context, id int64) (sql.NullInt64, error) {
	row := q.db.QueryRowContext(ctx, barExists, id)
	var column_1 sql.NullInt64
	err := row.Scan(&column_1)
	return column_1, err
}
