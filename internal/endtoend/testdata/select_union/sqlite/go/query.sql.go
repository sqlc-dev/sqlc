// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package querytest

import (
	"context"
)

const selectExcept = `-- name: SelectExcept :many
SELECT a, b FROM foo
EXCEPT
SELECT a, b FROM foo
`

func (q *Queries) SelectExcept(ctx context.Context) ([]Foo, error) {
	rows, err := q.db.QueryContext(ctx, selectExcept)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Foo
	for rows.Next() {
		var i Foo
		if err := rows.Scan(&i.A, &i.B); err != nil {
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

const selectIntersect = `-- name: SelectIntersect :many
SELECT a, b FROM foo
INTERSECT
SELECT a, b FROM foo
`

func (q *Queries) SelectIntersect(ctx context.Context) ([]Foo, error) {
	rows, err := q.db.QueryContext(ctx, selectIntersect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Foo
	for rows.Next() {
		var i Foo
		if err := rows.Scan(&i.A, &i.B); err != nil {
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

const selectUnion = `-- name: SelectUnion :many
SELECT a, b FROM foo
UNION
SELECT a, b FROM foo
`

func (q *Queries) SelectUnion(ctx context.Context) ([]Foo, error) {
	rows, err := q.db.QueryContext(ctx, selectUnion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Foo
	for rows.Next() {
		var i Foo
		if err := rows.Scan(&i.A, &i.B); err != nil {
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

const selectUnionOther = `-- name: SelectUnionOther :many
SELECT a, b FROM foo
UNION
SELECT a, b FROM bar
`

func (q *Queries) SelectUnionOther(ctx context.Context) ([]Foo, error) {
	rows, err := q.db.QueryContext(ctx, selectUnionOther)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Foo
	for rows.Next() {
		var i Foo
		if err := rows.Scan(&i.A, &i.B); err != nil {
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

const selectUnionWithLimit = `-- name: SelectUnionWithLimit :many
SELECT a, b FROM foo
UNION
SELECT a, b FROM foo
LIMIT ? OFFSET ?
`

type SelectUnionWithLimitParams struct {
	Limit  int64
	Offset int64
}

func (q *Queries) SelectUnionWithLimit(ctx context.Context, arg SelectUnionWithLimitParams) ([]Foo, error) {
	rows, err := q.db.QueryContext(ctx, selectUnionWithLimit, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Foo
	for rows.Next() {
		var i Foo
		if err := rows.Scan(&i.A, &i.B); err != nil {
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
