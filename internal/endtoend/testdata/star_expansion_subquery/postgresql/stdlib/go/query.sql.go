// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
)

const starExpansionSubquery = `-- name: StarExpansionSubquery :many
SELECT a, b FROM foo WHERE EXISTS (SELECT a, b FROM foo)
`

func (q *Queries) StarExpansionSubquery(ctx context.Context) ([]Foo, error) {
	ctx, done := q.observer(ctx, "StarExpansionSubquery")
	rows, err := q.db.QueryContext(ctx, starExpansionSubquery)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []Foo
	for rows.Next() {
		var i Foo
		if err := rows.Scan(&i.A, &i.B); err != nil {
			return nil, done(err)
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}
