// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"
)

const twoJoins = `-- name: TwoJoins :many
SELECT foo.bar_id, foo.baz_id
FROM foo
JOIN bar ON bar.id = bar_id
JOIN baz ON baz.id = baz_id
`

func (q *Queries) TwoJoins(ctx context.Context) ([]Foo, error) {
	rows, err := q.db.Query(ctx, twoJoins)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Foo
	for rows.Next() {
		var i Foo
		if err := rows.Scan(&i.BarID, &i.BazID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
