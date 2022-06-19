// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
)

const schemaScopedList = `-- name: SchemaScopedList :many
SELECT id FROM foo.bar
`

func (q *Queries) SchemaScopedList(ctx context.Context) ([]int32, error) {
	ctx, done := q.observer(ctx, "SchemaScopedList")
	rows, err := q.db.QueryContext(ctx, schemaScopedList)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, done(err)
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}
