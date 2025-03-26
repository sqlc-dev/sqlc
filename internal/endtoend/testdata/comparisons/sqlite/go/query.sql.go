// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package querytest

import (
	"context"
)

const alsoNotEqual = `-- name: AlsoNotEqual :many
SELECT count(*) <> 0 FROM bar
`

func (q *Queries) AlsoNotEqual(ctx context.Context) ([]bool, error) {
	rows, err := q.db.QueryContext(ctx, alsoNotEqual)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []bool
	for rows.Next() {
		var column_1 bool
		if err := rows.Scan(&column_1); err != nil {
			return nil, err
		}
		items = append(items, column_1)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const equal = `-- name: Equal :many
SELECT count(*) = 0 FROM bar
`

func (q *Queries) Equal(ctx context.Context) ([]bool, error) {
	rows, err := q.db.QueryContext(ctx, equal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []bool
	for rows.Next() {
		var column_1 bool
		if err := rows.Scan(&column_1); err != nil {
			return nil, err
		}
		items = append(items, column_1)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const greaterThan = `-- name: GreaterThan :many
SELECT count(*) > 0 FROM bar
`

func (q *Queries) GreaterThan(ctx context.Context) ([]bool, error) {
	rows, err := q.db.QueryContext(ctx, greaterThan)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []bool
	for rows.Next() {
		var column_1 bool
		if err := rows.Scan(&column_1); err != nil {
			return nil, err
		}
		items = append(items, column_1)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const greaterThanOrEqual = `-- name: GreaterThanOrEqual :many
SELECT count(*) >= 0 FROM bar
`

func (q *Queries) GreaterThanOrEqual(ctx context.Context) ([]bool, error) {
	rows, err := q.db.QueryContext(ctx, greaterThanOrEqual)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []bool
	for rows.Next() {
		var column_1 bool
		if err := rows.Scan(&column_1); err != nil {
			return nil, err
		}
		items = append(items, column_1)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const lessThan = `-- name: LessThan :many
SELECT count(*) < 0 FROM bar
`

func (q *Queries) LessThan(ctx context.Context) ([]bool, error) {
	rows, err := q.db.QueryContext(ctx, lessThan)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []bool
	for rows.Next() {
		var column_1 bool
		if err := rows.Scan(&column_1); err != nil {
			return nil, err
		}
		items = append(items, column_1)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const lessThanOrEqual = `-- name: LessThanOrEqual :many
SELECT count(*) <= 0 FROM bar
`

func (q *Queries) LessThanOrEqual(ctx context.Context) ([]bool, error) {
	rows, err := q.db.QueryContext(ctx, lessThanOrEqual)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []bool
	for rows.Next() {
		var column_1 bool
		if err := rows.Scan(&column_1); err != nil {
			return nil, err
		}
		items = append(items, column_1)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const notEqual = `-- name: NotEqual :many
SELECT count(*) != 0 FROM bar
`

func (q *Queries) NotEqual(ctx context.Context) ([]bool, error) {
	rows, err := q.db.QueryContext(ctx, notEqual)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []bool
	for rows.Next() {
		var column_1 bool
		if err := rows.Scan(&column_1); err != nil {
			return nil, err
		}
		items = append(items, column_1)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
