// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package override

import (
	"context"
)

const test = `-- name: Test :one
SELECT 1
`

func (q *Queries) Test(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, test)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}
