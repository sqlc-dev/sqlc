// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package querytest

import (
	"context"
)

const listBar = `-- name: ListBar :many
SELECT id FROM bar
`

func (q *Queries) ListBar(ctx context.Context) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, listBar)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]int32, 0)
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listFoo = `-- name: ListFoo :many
SELECT id, bar FROM foo
`

func (q *Queries) ListFoo(ctx context.Context) ([]Foo, error) {
	rows, err := q.db.QueryContext(ctx, listFoo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Foo, 0)
	for rows.Next() {
		var i Foo
		if err := rows.Scan(&i.ID, &i.Bar); err != nil {
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
