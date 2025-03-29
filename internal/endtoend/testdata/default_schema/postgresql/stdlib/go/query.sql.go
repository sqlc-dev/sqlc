// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: query.sql

package querytest

import (
	"context"
)

const defaultSchemaCreate = `-- name: DefaultSchemaCreate :one
INSERT INTO foo.bar (id, name) VALUES ($1, $2) RETURNING id, name
`

type DefaultSchemaCreateParams struct {
	ID   int32
	Name string
}

func (q *Queries) DefaultSchemaCreate(ctx context.Context, arg DefaultSchemaCreateParams) (Bar, error) {
	row := q.db.QueryRowContext(ctx, defaultSchemaCreate, arg.ID, arg.Name)
	var i Bar
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const defaultSchemaSelect = `-- name: DefaultSchemaSelect :one
SELECT id, name FROM foo.bar WHERE id = $1
`

func (q *Queries) DefaultSchemaSelect(ctx context.Context, id int32) (Bar, error) {
	row := q.db.QueryRowContext(ctx, defaultSchemaSelect, id)
	var i Bar
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const defaultSchemaSelectAll = `-- name: DefaultSchemaSelectAll :many
SELECT id, name FROM foo.bar
`

func (q *Queries) DefaultSchemaSelectAll(ctx context.Context) ([]Bar, error) {
	rows, err := q.db.QueryContext(ctx, defaultSchemaSelectAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Bar
	for rows.Next() {
		var i Bar
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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
