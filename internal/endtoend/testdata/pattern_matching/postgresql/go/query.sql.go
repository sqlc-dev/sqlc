// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const petsByName = `-- name: PetsByName :many
SELECT name FROM pet WHERE name LIKE $1
`

func (q *Queries) PetsByName(ctx context.Context, name sql.NullString) ([]sql.NullString, error) {
	rows, err := q.db.QueryContext(ctx, petsByName, name)
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
