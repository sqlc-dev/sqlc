// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package querytest

import (
	"context"
)

const listBar = `-- name: ListBar :many
SELECT id FROM (
  SELECT id FROM bar
) bar
`

// Lists all bars
//
//	SELECT id FROM (
//	  SELECT id FROM bar
//	) bar
func (q *Queries) ListBar(ctx context.Context) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, listBar)
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

const removeBar = `-- name: RemoveBar :exec
DELETE FROM bar WHERE id = $1
`

// RemoveBar
//
//	DELETE FROM bar WHERE id = $1
func (q *Queries) RemoveBar(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, removeBar, id)
	return err
}
