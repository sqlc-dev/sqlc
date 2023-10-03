// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
)

const schemaScopedDelete = `-- name: SchemaScopedDelete :exec
DELETE FROM foo.bar WHERE id = ?
`

func (q *Queries) SchemaScopedDelete(ctx context.Context, id uint64) error {
	query := schemaScopedDelete
	queryParams := []interface{}{id}

	_, err := q.db.ExecContext(ctx, query, queryParams...)
	return err
}
