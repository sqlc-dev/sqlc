// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"
)

const joinNoConstraints = `-- name: JoinNoConstraints :many
SELECT foo.barid
FROM foo
CROSS JOIN bar
WHERE bar.id = ? AND owner = ?
`

type JoinNoConstraintsParams struct {
	ID    uint64
	Owner string
}

func (q *Queries) JoinNoConstraints(ctx context.Context, arg JoinNoConstraintsParams) ([]uint64, error) {
	rows, err := q.db.QueryContext(ctx, joinNoConstraints, arg.ID, arg.Owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uint64
	for rows.Next() {
		var barid uint64
		if err := rows.Scan(&barid); err != nil {
			return nil, err
		}
		items = append(items, barid)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const joinParamWhereClause = `-- name: JoinParamWhereClause :many
SELECT foo.barid
FROM foo
JOIN bar ON bar.id = ?
WHERE owner = ?
`

type JoinParamWhereClauseParams struct {
	ID    uint64
	Owner string
}

func (q *Queries) JoinParamWhereClause(ctx context.Context, arg JoinParamWhereClauseParams) ([]uint64, error) {
	rows, err := q.db.QueryContext(ctx, joinParamWhereClause, arg.ID, arg.Owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uint64
	for rows.Next() {
		var barid uint64
		if err := rows.Scan(&barid); err != nil {
			return nil, err
		}
		items = append(items, barid)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const joinWhereClause = `-- name: JoinWhereClause :many
SELECT foo.barid
FROM foo
JOIN bar ON bar.id = barid
WHERE owner = ?
`

func (q *Queries) JoinWhereClause(ctx context.Context, owner string) ([]uint64, error) {
	rows, err := q.db.QueryContext(ctx, joinWhereClause, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uint64
	for rows.Next() {
		var barid uint64
		if err := rows.Scan(&barid); err != nil {
			return nil, err
		}
		items = append(items, barid)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
