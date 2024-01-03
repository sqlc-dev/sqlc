// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package datatype

import (
	"context"
	"net/netip"
)

const find = `-- name: Find :one
SELECT bar FROM foo WHERE baz = $1
`

func (q *Queries) Find(ctx context.Context, baz *netip.Prefix) (*netip.Addr, error) {
	row := q.db.QueryRow(ctx, find, baz)
	var bar *netip.Addr
	err := row.Scan(&bar)
	return bar, err
}

const list = `-- name: List :many
SELECT bar, baz FROM foo
`

func (q *Queries) List(ctx context.Context) ([]Foo, error) {
	rows, err := q.db.Query(ctx, list)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Foo
	for rows.Next() {
		var i Foo
		if err := rows.Scan(&i.Bar, &i.Baz); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
