// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
)

const deleteBarByID = `-- name: DeleteBarByID :execrows
DELETE FROM bar WHERE id = $1
`

func (q *Queries) DeleteBarByID(ctx context.Context, id int32) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteBarByID, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
