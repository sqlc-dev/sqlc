// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const sumOrder = `-- name: SumOrder :one
SELECT SUM(quantity) FROM orders
`

func (q *Queries) SumOrder(ctx context.Context) (pgtype.Numeric, error) {
	row := q.db.QueryRow(ctx, sumOrder)
	var sum pgtype.Numeric
	err := row.Scan(&sum)
	return sum, err
}
