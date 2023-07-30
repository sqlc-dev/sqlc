// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
	"strings"

	"github.com/lib/pq"
)

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio, country_code, titles
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, name, bio, country_code, titles
`

type CreateAuthorParams struct {
	Name        string
	Bio         sql.NullString
	CountryCode string
	Titles      []string
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error) {
	row := q.db.QueryRowContext(ctx, createAuthor,
		arg.Name,
		arg.Bio,
		arg.CountryCode,
		pq.Array(arg.Titles),
	)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Bio,
		&i.CountryCode,
		pq.Array(&i.Titles),
	)
	return i, err
}

const createAuthorOnlyTitles = `-- name: CreateAuthorOnlyTitles :one
INSERT INTO authors (name, titles) VALUES ($1, $2) RETURNING id, name, bio, country_code, titles
`

func (q *Queries) CreateAuthorOnlyTitles(ctx context.Context, name string, titles []string) (Author, error) {
	row := q.db.QueryRowContext(ctx, createAuthorOnlyTitles, name, pq.Array(titles))
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Bio,
		&i.CountryCode,
		pq.Array(&i.Titles),
	)
	return i, err
}

const deleteAuthor = `-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1
`

func (q *Queries) DeleteAuthor(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAuthor, id)
	return err
}

const deleteAuthors = `-- name: DeleteAuthors :exec
DELETE FROM authors
WHERE id IN ($2) AND name = $1
`

func (q *Queries) DeleteAuthors(ctx context.Context, name string, ids []int64) error {
	query := deleteAuthors
	var queryParams []interface{}
	queryParams = append(queryParams, name)
	if len(ids) > 0 {
		for _, v := range ids {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:ids*/?", strings.Repeat(",?", len(ids))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:ids*/?", "NULL", 1)
	}
	_, err := q.db.ExecContext(ctx, query, queryParams...)
	return err
}

const getAuthor = `-- name: GetAuthor :one
SELECT id, name, bio, country_code, titles FROM authors
WHERE name = $1 AND country_code = $2 LIMIT 1
`

func (q *Queries) GetAuthor(ctx context.Context, name string, countryCode string) (Author, error) {
	row := q.db.QueryRowContext(ctx, getAuthor, name, countryCode)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Bio,
		&i.CountryCode,
		pq.Array(&i.Titles),
	)
	return i, err
}

const listAuthors = `-- name: ListAuthors :many
SELECT id, name, bio, country_code, titles FROM authors
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
			pq.Array(&i.Titles),
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
