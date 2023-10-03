// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
)

const percentile = `-- name: Percentile :one
select percentile_disc(0.5) within group (order by authors.name)
from authors
`

func (q *Queries) Percentile(ctx context.Context, aq ...AdditionalQuery) (interface{}, error) {
	query := percentile
	queryParams := []interface{}{}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var percentile_disc interface{}
	err := row.Scan(&percentile_disc)
	return percentile_disc, err
}
