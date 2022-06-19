// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
)

const schemaScopedFilter = `-- name: SchemaScopedFilter :many
SELECT id FROM foo.bar WHERE id = $1
`

func (q *Queries) SchemaScopedFilter(ctx context.Context, id int32) ([]int32, error) {
	ctx, done := q.observer(ctx, "SchemaScopedFilter")
	rows, err := q.db.Query(ctx, schemaScopedFilter, id)
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
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}
