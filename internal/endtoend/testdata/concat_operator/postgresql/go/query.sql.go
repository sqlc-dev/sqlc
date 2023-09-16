// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: query.sql

package querytest

import (
	"context"
)

const concat = `-- name: Concat :many
SELECT bar || ' ' || baz FROM foo
`

func (q *Queries) Concat(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, concat)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var column_1 string
		if err := rows.Scan(&column_1); err != nil {
			return nil, err
		}
		items = append(items, column_1)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
