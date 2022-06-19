// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
)

const placeholder = `-- name: Placeholder :many
SELECT name, location, size from venues
`

func (q *Queries) Placeholder(ctx context.Context) ([]Venue, error) {
	ctx, done := q.observer(ctx, "Placeholder")
	rows, err := q.db.QueryContext(ctx, placeholder)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []Venue
	for rows.Next() {
		var i Venue
		if err := rows.Scan(&i.Name, &i.Location, &i.Size); err != nil {
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
