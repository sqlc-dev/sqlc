// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package querytest

import (
	"context"
)

const getAuthorByID = `-- name: GetAuthorByID :one
SELECT  id, name, bio
FROM    authors
WHERE   id = ?
LIMIT   1
`

func (q *Queries) GetAuthorByID(ctx context.Context, id int64) (Author, error) {
	row := q.db.QueryRowContext(ctx, getAuthorByID, id)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return i, err
}

const getAuthorIDByID = `-- name: GetAuthorIDByID :one
SELECT  id
FROM    authors
WHERE   id = ?
LIMIT   1
`

func (q *Queries) GetAuthorIDByID(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getAuthorIDByID, id)
	err := row.Scan(&id)
	return id, err
}

const getUser = `-- name: GetUser :one
SELECT  sub
FROM    users
WHERE   sub = ?
LIMIT   1
`

func (q *Queries) GetUser(ctx context.Context, sub string) (string, error) {
	row := q.db.QueryRowContext(ctx, getUser, sub)
	err := row.Scan(&sub)
	return sub, err
}
