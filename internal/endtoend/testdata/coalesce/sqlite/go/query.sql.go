// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const coalesce = `-- name: Coalesce :many
SELECT coalesce(bar, '') as login
FROM foo
`

func (q *Queries) Coalesce(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, coalesce)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var login string
		if err := rows.Scan(&login); err != nil {
			return nil, err
		}
		items = append(items, login)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const coalesceColumns = `-- name: CoalesceColumns :many
SELECT bar, bat, coalesce(bar, bat)
FROM foo
`

type CoalesceColumnsRow struct {
	Bar   sql.NullString
	Bat   string
	Bar_2 string
}

func (q *Queries) CoalesceColumns(ctx context.Context) ([]CoalesceColumnsRow, error) {
	rows, err := q.db.QueryContext(ctx, coalesceColumns)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CoalesceColumnsRow
	for rows.Next() {
		var i CoalesceColumnsRow
		if err := rows.Scan(&i.Bar, &i.Bat, &i.Bar_2); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
