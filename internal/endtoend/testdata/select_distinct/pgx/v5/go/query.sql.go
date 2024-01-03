// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package querytest

import (
	"context"
)

const getBars = `-- name: GetBars :many
SELECT DISTINCT ON (a.id) a.id, a.name
FROM bar a
`

func (q *Queries) GetBars(ctx context.Context) ([]Bar, error) {
	rows, err := q.db.Query(ctx, getBars)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Bar
	for rows.Next() {
		var i Bar
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
