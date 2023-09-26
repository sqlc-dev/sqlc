// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package db

import (
	"context"
	"database/sql"
)

const upsertAuthor = `-- name: UpsertAuthor :exec
INSERT INTO authors (name, bio)
VALUES ($1, $2)
ON CONFLICT (name) DO UPDATE
SET bio = $2
`

type UpsertAuthorParams struct {
	Name string
	Bio  sql.NullString
}

func (q *Queries) UpsertAuthor(ctx context.Context, arg UpsertAuthorParams) error {
	_, err := q.db.ExecContext(ctx, upsertAuthor, arg.Name, arg.Bio)
	return err
}

const upsertAuthorNamed = `-- name: UpsertAuthorNamed :exec
INSERT INTO authors (name, bio)
VALUES ($1, $2)
ON CONFLICT (name) DO UPDATE
SET bio = $2
`

type UpsertAuthorNamedParams struct {
	Name string
	Bio  sql.NullString
}

func (q *Queries) UpsertAuthorNamed(ctx context.Context, arg UpsertAuthorNamedParams) error {
	_, err := q.db.ExecContext(ctx, upsertAuthorNamed, arg.Name, arg.Bio)
	return err
}
