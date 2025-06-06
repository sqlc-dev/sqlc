// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"
)

const tableName = `-- name: TableName :one
SELECT foo.id
FROM foo
JOIN bar ON foo.bar = bar.id
WHERE bar.id = $1 AND foo.id = $2
`

type TableNameParams struct {
	ID   int32
	ID_2 int32
}

func (q *Queries) TableName(ctx context.Context, arg TableNameParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, tableName, arg.ID, arg.ID_2)
	var id int32
	err := row.Scan(&id)
	return id, err
}
