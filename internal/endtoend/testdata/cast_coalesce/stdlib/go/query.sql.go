// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
)

const castCoalesce = `-- name: CastCoalesce :many
SELECT coalesce(bar, '')::text as login
FROM foo
`

func (q *Queries) CastCoalesce(ctx context.Context) ([]string, error) {
	ctx, done := q.observer(ctx, "CastCoalesce")
	rows, err := q.db.QueryContext(ctx, castCoalesce)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var login string
		if err := rows.Scan(&login); err != nil {
			return nil, done(err)
		}
		items = append(items, login)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}
