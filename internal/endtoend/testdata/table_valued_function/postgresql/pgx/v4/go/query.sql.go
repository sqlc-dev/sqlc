// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const getAuthorNames = `-- name: GetAuthorNames :many
SELECT a.name  FROM fauthors() a
`

func (q *Queries) GetAuthorNames(ctx context.Context) ([]sql.NullString, error) {
	rows, err := q.db.Query(ctx, getAuthorNames)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
