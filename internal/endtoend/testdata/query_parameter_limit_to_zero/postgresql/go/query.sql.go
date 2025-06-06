// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio, country_code
) VALUES (
  $1, $2, $3
)
RETURNING id, name, bio, country_code
`

type CreateAuthorParams struct {
	Name        string
	Bio         sql.NullString
	CountryCode string
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error) {
	row := q.db.QueryRowContext(ctx, createAuthor, arg.Name, arg.Bio, arg.CountryCode)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Bio,
		&i.CountryCode,
	)
	return i, err
}

const deleteAuthor = `-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1
`

type DeleteAuthorParams struct {
	ID int64
}

func (q *Queries) DeleteAuthor(ctx context.Context, arg DeleteAuthorParams) error {
	_, err := q.db.ExecContext(ctx, deleteAuthor, arg.ID)
	return err
}

const getAuthor = `-- name: GetAuthor :one
SELECT id, name, bio, country_code FROM authors
WHERE name = $1 AND country_code = $2 LIMIT 1
`

type GetAuthorParams struct {
	Name        string
	CountryCode string
}

func (q *Queries) GetAuthor(ctx context.Context, arg GetAuthorParams) (Author, error) {
	row := q.db.QueryRowContext(ctx, getAuthor, arg.Name, arg.CountryCode)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Bio,
		&i.CountryCode,
	)
	return i, err
}

const listAuthors = `-- name: ListAuthors :many
SELECT id, name, bio, country_code FROM authors
ORDER BY name
`

func (q *Queries) ListAuthors(ctx context.Context) ([]Author, error) {
	rows, err := q.db.QueryContext(ctx, listAuthors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Author
	for rows.Next() {
		var i Author
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Bio,
			&i.CountryCode,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
