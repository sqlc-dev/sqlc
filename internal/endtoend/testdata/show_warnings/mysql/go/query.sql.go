// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package querytest

import (
	"context"
)

const showWarnings = `-- name: ShowWarnings :many
SHOW WARNINGS
`

type ShowWarningsRow struct {
	Level   string
	Code    int32
	Message string
}

func (q *Queries) ShowWarnings(ctx context.Context) ([]ShowWarningsRow, error) {
	rows, err := q.db.QueryContext(ctx, showWarnings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ShowWarningsRow
	for rows.Next() {
		var i ShowWarningsRow
		if err := rows.Scan(&i.Level, &i.Code, &i.Message); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
