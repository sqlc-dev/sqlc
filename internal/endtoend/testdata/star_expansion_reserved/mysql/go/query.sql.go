// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
)

const starExpansionReserved = `-- name: StarExpansionReserved :many
SELECT ` + "`" + `group` + "`" + `, ` + "`" + `key` + "`" + ` FROM foo
`

func (q *Queries) StarExpansionReserved(ctx context.Context) ([]Foo, error) {
	ctx, done := q.observer(ctx, "StarExpansionReserved")
	rows, err := q.db.QueryContext(ctx, starExpansionReserved)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []Foo
	for rows.Next() {
		var i Foo
		if err := rows.Scan(&i.Group, &i.Key); err != nil {
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
