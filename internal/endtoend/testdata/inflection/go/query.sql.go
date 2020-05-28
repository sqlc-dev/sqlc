// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package querytest

import (
	"context"
)

const listCampus = `-- name: ListCampus :many
SELECT id FROM campus
`

func (q *Queries) ListCampus(ctx context.Context) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, listCampus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]int32, 0)
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
