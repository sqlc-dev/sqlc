// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.1
// source: query.sql

package querytest

import (
	"context"
	"time"

	"github.com/jackc/pgtype"
)

const generateSeries = `-- name: GenerateSeries :many
SELECT generate_series($1::timestamp, $2::timestamp)
`

type GenerateSeriesParams struct {
	Column1 time.Time `json:"column_1"`
	Column2 time.Time `json:"column_2"`
}

func (q *Queries) GenerateSeries(ctx context.Context, arg GenerateSeriesParams) ([]pgtype.Numeric, error) {
	rows, err := q.db.Query(ctx, generateSeries, arg.Column1, arg.Column2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.Numeric
	for rows.Next() {
		var generate_series pgtype.Numeric
		if err := rows.Scan(&generate_series); err != nil {
			return nil, err
		}
		items = append(items, generate_series)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
