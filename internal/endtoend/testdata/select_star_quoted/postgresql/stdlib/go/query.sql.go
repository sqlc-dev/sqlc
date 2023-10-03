// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const getAll = `-- name: GetAll :many
SELECT "CamelCase" FROM users
`

func (q *Queries) GetAll(ctx context.Context, aq ...AdditionalQuery) ([]sql.NullString, error) {
	query := getAll
	queryParams := []interface{}{}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var CamelCase sql.NullString
		if err := rows.Scan(&CamelCase); err != nil {
			return nil, err
		}
		items = append(items, CamelCase)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
