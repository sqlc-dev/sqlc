// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package querytest

import (
	"context"
)

const advisoryLock = `-- name: AdvisoryLock :many
SELECT pg_advisory_unlock($1)
`

func (q *Queries) AdvisoryLock(ctx context.Context, key int64) ([]bool, error) {
	rows, err := q.db.QueryContext(ctx, advisoryLock, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]bool, 0)
	for rows.Next() {
		var pg_advisory_unlock bool
		if err := rows.Scan(&pg_advisory_unlock); err != nil {
			return nil, err
		}
		items = append(items, pg_advisory_unlock)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
