// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const updateJ = `-- name: UpdateJ :exec
UPDATE
    test
SET
    j = jsonb_build_object($1, $2, $3, $4)
WHERE
    id = $5
`

type UpdateJParams struct {
	JsonbBuildObject   interface{}
	JsonbBuildObject_2 interface{}
	JsonbBuildObject_3 interface{}
	JsonbBuildObject_4 interface{}
	ID                 sql.NullInt32
}

func (q *Queries) UpdateJ(ctx context.Context, arg UpdateJParams) error {
	query := updateJ
	queryParams := []interface{}{
		arg.JsonbBuildObject,
		arg.JsonbBuildObject_2,
		arg.JsonbBuildObject_3,
		arg.JsonbBuildObject_4,
		arg.ID,
	}

	_, err := q.db.ExecContext(ctx, query, queryParams...)
	return err
}
