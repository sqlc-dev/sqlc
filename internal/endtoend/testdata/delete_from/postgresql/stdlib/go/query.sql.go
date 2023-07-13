// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: query.sql

package querytest

import (
	"context"
)

const deleteFrom = `-- name: DeleteFrom :exec
DELETE FROM foo WHERE id = $1
`

func (q *Queries) DeleteFrom(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteFrom, id)
	return err
}
