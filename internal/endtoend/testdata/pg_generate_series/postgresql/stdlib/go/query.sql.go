// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: query.sql

package querytest

import (
	"context"
	"time"
)

const generateSeries = `-- name: GenerateSeries :many
SELECT generate_series($1::timestamp, $2::timestamp)
`

type GenerateSeriesParams struct {
	Column1 time.Time `json:"column_1"`
	Column2 time.Time `json:"column_2"`
}

func (q *Queries) GenerateSeries(ctx context.Context, arg GenerateSeriesParams) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, generateSeries, arg.Column1, arg.Column2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var generate_series string
		if err := rows.Scan(&generate_series); err != nil {
			return nil, err
		}
		items = append(items, generate_series)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
