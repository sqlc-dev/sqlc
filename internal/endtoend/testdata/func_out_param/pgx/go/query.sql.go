// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAuthor = `-- name: CreateAuthor :one
SELECT id FROM add_author (
  $1, $2
)
`

type CreateAuthorParams struct {
	Name string
	Bio  string
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (pgtype.Int4, error) {
	row := q.db.QueryRow(ctx, createAuthor, arg.Name, arg.Bio)
	var id pgtype.Int4
	err := row.Scan(&id)
	return id, err
}
