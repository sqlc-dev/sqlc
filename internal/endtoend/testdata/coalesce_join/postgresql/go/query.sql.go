// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
)

const getBar = `-- name: GetBar :many
SELECT foo.id, COALESCE(bar.id, 0) AS bar_id
FROM foo
LEFT JOIN bar ON foo.id = bar.id
`

type GetBarRow struct {
	ID    int64
	BarID int64
}

func (q *Queries) GetBar(ctx context.Context) ([]GetBarRow, error) {
	rows, err := q.db.QueryContext(ctx, getBar)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBarRow
	for rows.Next() {
		var i GetBarRow
		if err := rows.Scan(&i.ID, &i.BarID); err != nil {
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
