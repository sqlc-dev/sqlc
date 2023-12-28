// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getDate = `-- name: GetDate :one
SELECT now from NOW()
`

func (q *Queries) GetDate(ctx context.Context) (pgtype.Timestamptz, error) {
	row := q.db.QueryRow(ctx, getDate)
	var now pgtype.Timestamptz
	err := row.Scan(&now)
	return now, err
}
