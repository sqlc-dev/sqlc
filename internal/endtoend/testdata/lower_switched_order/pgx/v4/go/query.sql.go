// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"
)

const lowerSwitchedOrder = `-- name: LowerSwitchedOrder :many
SELECT bar FROM foo WHERE bar = $1 AND bat = LOWER($2)
`

type LowerSwitchedOrderParams struct {
	Bar   string
	Lower string
}

func (q *Queries) LowerSwitchedOrder(ctx context.Context, arg LowerSwitchedOrderParams) ([]string, error) {
	rows, err := q.db.Query(ctx, lowerSwitchedOrder, arg.Bar, arg.Lower)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var bar string
		if err := rows.Scan(&bar); err != nil {
			return nil, err
		}
		items = append(items, bar)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
