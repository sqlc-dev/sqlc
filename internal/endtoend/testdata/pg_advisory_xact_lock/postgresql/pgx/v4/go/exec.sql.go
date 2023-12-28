// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: exec.sql

package querytest

import (
	"context"
)

const advisoryLockExec = `-- name: AdvisoryLockExec :exec
SELECT pg_advisory_lock($1)
`

func (q *Queries) AdvisoryLockExec(ctx context.Context, pgAdvisoryLock int64) error {
	_, err := q.db.Exec(ctx, advisoryLockExec, pgAdvisoryLock)
	return err
}

const advisoryLockExecRows = `-- name: AdvisoryLockExecRows :execrows
SELECT pg_advisory_lock($1)
`

func (q *Queries) AdvisoryLockExecRows(ctx context.Context, pgAdvisoryLock int64) (int64, error) {
	result, err := q.db.Exec(ctx, advisoryLockExecRows, pgAdvisoryLock)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
