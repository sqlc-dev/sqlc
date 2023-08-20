// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: query.sql

package authors

import (
	"context"
	"database/sql"
)

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO authors (
          name, bio
) VALUES (
  $1, $2
)
RETURNING id, name, bio
`

type CreateAuthorParams struct {
	Name string
	Bio  sql.NullString
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams, aq ...AdditionalQuery) (Author, error) {
	query := createAuthor
	queryParams := []interface{}{arg.Name, arg.Bio}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return i, err
}

const getAuthor = `-- name: GetAuthor :one
SELECT id, name, bio FROM authors
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAuthor(ctx context.Context, id int64, aq ...AdditionalQuery) (Author, error) {
	query := getAuthor
	queryParams := []interface{}{id}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return i, err
}

const listAuthors = `-- name: ListAuthors :many
SELECT id, name, bio FROM authors
ORDER BY bio
`

func (q *Queries) ListAuthors(ctx context.Context, aq ...AdditionalQuery) ([]Author, error) {
	query := listAuthors
	queryParams := []interface{}{}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Author
	for rows.Next() {
		var i Author
		if err := rows.Scan(&i.ID, &i.Name, &i.Bio); err != nil {
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

const selectOne = `-- name: SelectOne :one
SELECT 1
`

func (q *Queries) SelectOne(ctx context.Context, aq ...AdditionalQuery) (int32, error) {
	query := selectOne
	queryParams := []interface{}{}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var column_1 int32
	err := row.Scan(&column_1)
	return column_1, err
}
