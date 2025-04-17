// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"
)

const in = `-- name: In :many
SELECT id
FROM bar
WHERE id IN ($1, $2)
`

type InParams struct {
	ID   int32
	ID_2 int32
}

func (q *Queries) In(ctx context.Context, arg InParams) ([]int32, error) {
	rows, err := q.db.Query(ctx, in, arg.ID, arg.ID_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
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
