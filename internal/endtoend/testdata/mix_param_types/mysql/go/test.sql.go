// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: test.sql

package querytest

import (
	"context"
)

const countOne = `-- name: CountOne :one
SELECT count(1) FROM bar WHERE id = ? AND name <> ?
`

type CountOneParams struct {
	ID   int64
	Name string
}

func (q *Queries) CountOne(ctx context.Context, arg CountOneParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, countOne, arg.ID, arg.Name)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countThree = `-- name: CountThree :one
SELECT count(1) FROM bar WHERE id > ? AND phone <> ? AND name <> ?
`

type CountThreeParams struct {
	ID    int64
	Phone string
	Name  string
}

func (q *Queries) CountThree(ctx context.Context, arg CountThreeParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, countThree, arg.ID, arg.Phone, arg.Name)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countTwo = `-- name: CountTwo :one
SELECT count(1) FROM bar WHERE id = ? AND name <> ?
`

type CountTwoParams struct {
	ID   int64
	Name string
}

func (q *Queries) CountTwo(ctx context.Context, arg CountTwoParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, countTwo, arg.ID, arg.Name)
	var count int64
	err := row.Scan(&count)
	return count, err
}
