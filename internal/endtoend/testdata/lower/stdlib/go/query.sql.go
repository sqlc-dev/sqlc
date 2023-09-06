// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: query.sql

package querytest

import (
	"context"
)

const lower = `-- name: Lower :many
SELECT bar FROM foo WHERE bar = $1 AND LOWER(bat) = $2
`

type LowerParams struct {
	Bar string
	Bat string
}

func (q *Queries) Lower(ctx context.Context, arg LowerParams) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, lower, arg.Bar, arg.Bat)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var bar string
		if err := rows.Scan(&bar); err != nil {
			return nil, err
		}
		items = append(items, bar)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
