// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
)

const textArray = `-- name: TextArray :many
SELECT tags FROM bar
`

func (q *Queries) TextArray(ctx context.Context, aq ...AdditionalQuery) ([][][]string, error) {
	query := textArray
	queryParams := []interface{}{}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	rows, err := q.db.Query(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items [][][]string
	for rows.Next() {
		var tags [][]string
		if err := rows.Scan(&tags); err != nil {
			return nil, err
		}
		items = append(items, tags)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
