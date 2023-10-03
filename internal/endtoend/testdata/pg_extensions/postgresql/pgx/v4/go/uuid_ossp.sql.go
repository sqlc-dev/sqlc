// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: uuid_ossp.sql

package querytest

import (
	"context"

	"github.com/google/uuid"
)

const generateUUID = `-- name: GenerateUUID :one
SELECT uuid_generate_v4()
`

func (q *Queries) GenerateUUID(ctx context.Context, aq ...AdditionalQuery) (uuid.UUID, error) {
	query := generateUUID
	queryParams := []interface{}{}
	row := q.db.QueryRow(ctx, query, queryParams...)
	var uuid_generate_v4 uuid.UUID
	err := row.Scan(&uuid_generate_v4)
	return uuid_generate_v4, err
}
