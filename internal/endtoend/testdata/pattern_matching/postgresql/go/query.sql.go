// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const petsByName = `-- name: PetsByName :many
SELECT name FROM pet WHERE name LIKE $1
`

func (q *Queries) PetsByName(ctx context.Context, name sql.NullString, aq ...AdditionalQuery) ([]sql.NullString, error) {
	query := petsByName
	queryParams := []interface{}{name}

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
		var name sql.NullString
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
