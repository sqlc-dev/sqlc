// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const identOnNonNullable = `-- name: IdentOnNonNullable :many
SELECT bar FROM foo WHERE bar = $1
`

func (q *Queries) IdentOnNonNullable(ctx context.Context, bar pgtype.Text) ([]string, error) {
	rows, err := q.db.Query(ctx, identOnNonNullable, bar)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var bar string
		if err := rows.Scan(&bar); err != nil {
			return nil, err
		}
		items = append(items, bar)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const identOnNullable = `-- name: IdentOnNullable :many
SELECT maybe_bar FROM foo WHERE maybe_bar = $1
`

func (q *Queries) IdentOnNullable(ctx context.Context, maybeBar pgtype.Text) ([]pgtype.Text, error) {
	rows, err := q.db.Query(ctx, identOnNullable, maybeBar)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.Text
	for rows.Next() {
		var maybe_bar pgtype.Text
		if err := rows.Scan(&maybe_bar); err != nil {
			return nil, err
		}
		items = append(items, maybe_bar)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const stringOnNonNullable = `-- name: StringOnNonNullable :many
SELECT bar FROM foo WHERE bar = $1
`

func (q *Queries) StringOnNonNullable(ctx context.Context, bar pgtype.Text) ([]string, error) {
	rows, err := q.db.Query(ctx, stringOnNonNullable, bar)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var bar string
		if err := rows.Scan(&bar); err != nil {
			return nil, err
		}
		items = append(items, bar)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const stringOnNullable = `-- name: StringOnNullable :many
SELECT maybe_bar FROM foo WHERE maybe_bar = $1
`

func (q *Queries) StringOnNullable(ctx context.Context, maybeBar pgtype.Text) ([]pgtype.Text, error) {
	rows, err := q.db.Query(ctx, stringOnNullable, maybeBar)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.Text
	for rows.Next() {
		var maybe_bar pgtype.Text
		if err := rows.Scan(&maybe_bar); err != nil {
			return nil, err
		}
		items = append(items, maybe_bar)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
