// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package querytest

import (
	"context"
)

const as = `-- name: As :many
SELECT name, name AS "other_name" FROM foo
`

type AsRow struct {
	Name      string
	OtherName string
}

func (q *Queries) As(ctx context.Context) ([]AsRow, error) {
	rows, err := q.db.QueryContext(ctx, as)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]AsRow, 0)
	for rows.Next() {
		var i AsRow
		if err := rows.Scan(&i.Name, &i.OtherName); err != nil {
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
