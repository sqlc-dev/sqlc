// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
)

const listAuthors = `-- name: ListAuthors :many
SELECT id, status, level FROM log_lines
`

func (q *Queries) ListAuthors(ctx context.Context, aq ...AdditionalQuery) ([]LogLine, error) {
	query := listAuthors
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
	var items []LogLine
	for rows.Next() {
		var i LogLine
		if err := rows.Scan(&i.ID, &i.Status, &i.Level); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
