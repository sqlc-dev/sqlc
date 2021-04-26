// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package querytest

import (
	"context"
)

const getMayors = `-- name: GetMayors :many
SELECT
    user_id,
    mayors.full_name
FROM users
LEFT JOIN cities USING (city_id)
INNER JOIN mayors USING (mayor_id)
`

type GetMayorsRow struct {
	UserID   int32
	FullName string
}

func (q *Queries) GetMayors(ctx context.Context) ([]GetMayorsRow, error) {
	rows, err := q.db.QueryContext(ctx, getMayors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMayorsRow
	for rows.Next() {
		var i GetMayorsRow
		if err := rows.Scan(&i.UserID, &i.FullName); err != nil {
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
