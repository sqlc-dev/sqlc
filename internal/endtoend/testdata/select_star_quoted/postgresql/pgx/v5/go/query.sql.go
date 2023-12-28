// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getAll = `-- name: GetAll :many
SELECT "CamelCase" FROM users
`

func (q *Queries) GetAll(ctx context.Context) ([]pgtype.Text, error) {
	rows, err := q.db.Query(ctx, getAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.Text
	for rows.Next() {
		var CamelCase pgtype.Text
		if err := rows.Scan(&CamelCase); err != nil {
			return nil, err
		}
		items = append(items, CamelCase)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
