// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
)

const listBar = `-- name: ListBar :many
SELECT id FROM bar
`

func (q *Queries) ListBar(ctx context.Context, aq ...AdditionalQuery) ([]int32, error) {
	query := listBar
	queryParams := []interface{}{}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	rows, err := q.db.Query(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int32{}
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
