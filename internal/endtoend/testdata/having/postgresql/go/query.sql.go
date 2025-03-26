// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package querytest

import (
	"context"
)

const coldCities = `-- name: ColdCities :many
SELECT city
FROM weather
GROUP BY city
HAVING max(temp_lo) < $1
`

func (q *Queries) ColdCities(ctx context.Context, tempLo int32) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, coldCities, tempLo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var city string
		if err := rows.Scan(&city); err != nil {
			return nil, err
		}
		items = append(items, city)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
