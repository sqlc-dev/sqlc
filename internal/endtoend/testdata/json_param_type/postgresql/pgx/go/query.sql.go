// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package querytest

import (
	"context"
)

const findByAddress = `-- name: FindByAddress :one
SELECT id, metadata FROM "user" WHERE "metadata"->>'address1' = $1 LIMIT 1
`

func (q *Queries) FindByAddress(ctx context.Context, metadata []byte) (User, error) {
	row := q.db.QueryRow(ctx, findByAddress, metadata)
	var i User
	err := row.Scan(&i.ID, &i.Metadata)
	return i, err
}
