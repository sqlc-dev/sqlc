// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const identOnNonNullable = `-- name: IdentOnNonNullable :many
SELECT bar FROM foo WHERE bar = ?
`

func (q *Queries) IdentOnNonNullable(ctx context.Context, bar sql.NullString) ([]string, error) {
	ctx, done := q.observer(ctx, "IdentOnNonNullable")
	rows, err := q.db.QueryContext(ctx, identOnNonNullable, bar)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var bar string
		if err := rows.Scan(&bar); err != nil {
			return nil, done(err)
		}
		items = append(items, bar)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}

const identOnNullable = `-- name: IdentOnNullable :many
SELECT maybe_bar FROM foo WHERE maybe_bar = ?
`

func (q *Queries) IdentOnNullable(ctx context.Context, maybeBar sql.NullString) ([]sql.NullString, error) {
	ctx, done := q.observer(ctx, "IdentOnNullable")
	rows, err := q.db.QueryContext(ctx, identOnNullable, maybeBar)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var maybe_bar sql.NullString
		if err := rows.Scan(&maybe_bar); err != nil {
			return nil, done(err)
		}
		items = append(items, maybe_bar)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}

const stringOnNonNullable = `-- name: StringOnNonNullable :many
SELECT bar FROM foo WHERE bar = ?
`

func (q *Queries) StringOnNonNullable(ctx context.Context, bar sql.NullString) ([]string, error) {
	ctx, done := q.observer(ctx, "StringOnNonNullable")
	rows, err := q.db.QueryContext(ctx, stringOnNonNullable, bar)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var bar string
		if err := rows.Scan(&bar); err != nil {
			return nil, done(err)
		}
		items = append(items, bar)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}

const stringOnNullable = `-- name: StringOnNullable :many
SELECT maybe_bar FROM foo WHERE maybe_bar = ?
`

func (q *Queries) StringOnNullable(ctx context.Context, maybeBar sql.NullString) ([]sql.NullString, error) {
	ctx, done := q.observer(ctx, "StringOnNullable")
	rows, err := q.db.QueryContext(ctx, stringOnNullable, maybeBar)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var maybe_bar sql.NullString
		if err := rows.Scan(&maybe_bar); err != nil {
			return nil, done(err)
		}
		items = append(items, maybe_bar)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}
