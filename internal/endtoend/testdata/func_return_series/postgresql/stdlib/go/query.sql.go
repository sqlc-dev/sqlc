// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
)

const generateSeries = `-- name: GenerateSeries :many
SELECT ($1::int) + i
FROM generate_series(0, $2::int) AS i
LIMIT 1
`

type GenerateSeriesParams struct {
	Column1 int32
	Column2 int32
}

func (q *Queries) GenerateSeries(ctx context.Context, arg GenerateSeriesParams) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, generateSeries, arg.Column1, arg.Column2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var column_1 int32
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
