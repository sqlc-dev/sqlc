// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
)

const barNotExists = `-- name: BarNotExists :one
SELECT
    NOT EXISTS (
        SELECT
            1
        FROM
            bar
        where
            id = $1
    )
`

func (q *Queries) BarNotExists(ctx context.Context, id int32, aq ...AdditionalQuery) (bool, error) {
	query := barNotExists
	queryParams := []interface{}{id}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var not_exists bool
	err := row.Scan(&not_exists)
	return not_exists, err
}
