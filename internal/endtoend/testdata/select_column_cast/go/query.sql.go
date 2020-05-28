// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package querytest

import (
	"context"
)

const selectColumnCast = `-- name: SelectColumnCast :many
SELECT bar::int FROM foo
`

func (q *Queries) SelectColumnCast(ctx context.Context) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, selectColumnCast)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]int32, 0)
	for rows.Next() {
		var bar int32
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
