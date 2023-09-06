// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: query.sql

package querytest

import (
	"context"
)

const limitMe = `-- name: LimitMe :many
SELECT bar FROM foo LIMIT $1
`

func (q *Queries) LimitMe(ctx context.Context, limit int32) ([]bool, error) {
	rows, err := q.db.QueryContext(ctx, limitMe, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []bool
	for rows.Next() {
		var bar bool
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
