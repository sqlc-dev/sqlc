// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const advisoryLockExecResult = `-- name: AdvisoryLockExecResult :execresult
SELECT pg_advisory_lock($1)
`

func (q *Queries) AdvisoryLockExecResult(ctx context.Context, pgAdvisoryLock int64) (sql.Result, error) {
	return q.db.ExecContext(ctx, advisoryLockExecResult, pgAdvisoryLock)
}

const advisoryLockOne = `-- name: AdvisoryLockOne :one
SELECT pg_advisory_lock($1)
`

func (q *Queries) AdvisoryLockOne(ctx context.Context, pgAdvisoryLock int64) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, advisoryLockOne, pgAdvisoryLock)
	var pg_advisory_lock interface{}
	err := row.Scan(&pg_advisory_lock)
	return pg_advisory_lock, err
}

const advisoryUnlock = `-- name: AdvisoryUnlock :many
SELECT pg_advisory_unlock($1)
`

func (q *Queries) AdvisoryUnlock(ctx context.Context, pgAdvisoryUnlock int64) ([]bool, error) {
	rows, err := q.db.QueryContext(ctx, advisoryUnlock, pgAdvisoryUnlock)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []bool
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
