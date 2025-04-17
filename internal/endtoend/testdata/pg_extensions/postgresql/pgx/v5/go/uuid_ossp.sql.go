// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: uuid_ossp.sql

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const generateUUID = `-- name: GenerateUUID :one
SELECT uuid_generate_v4()
`

func (q *Queries) GenerateUUID(ctx context.Context) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, generateUUID)
	var uuid_generate_v4 pgtype.UUID
	err := row.Scan(&uuid_generate_v4)
	return uuid_generate_v4, err
}
