// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const doubleDash = `-- name: DoubleDash :one
SELECT bar FROM foo LIMIT 1
`

func (q *Queries) DoubleDash(ctx context.Context, aq ...AdditionalQuery) (sql.NullString, error) {
	query := doubleDash
	queryParams := []interface{}{}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var bar sql.NullString
	err := row.Scan(&bar)
	return bar, err
}

const hash = `-- name: Hash :one
# name: Hash :one
SELECT bar FROM foo LIMIT 1
`

func (q *Queries) Hash(ctx context.Context, aq ...AdditionalQuery) (sql.NullString, error) {
	query := hash
	queryParams := []interface{}{}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var bar sql.NullString
	err := row.Scan(&bar)
	return bar, err
}

const slashStar = `-- name: SlashStar :one
SELECT bar FROM foo LIMIT 1
`

func (q *Queries) SlashStar(ctx context.Context, aq ...AdditionalQuery) (sql.NullString, error) {
	query := slashStar
	queryParams := []interface{}{}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var bar sql.NullString
	err := row.Scan(&bar)
	return bar, err
}
