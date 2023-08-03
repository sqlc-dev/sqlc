// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countRowsEmptyTable = `-- name: CountRowsEmptyTable :many
SELECT a, (SELECT count(a) FROM empty) as "count" FROM foo
`

type CountRowsEmptyTableRow struct {
	A     int32
	Count int64
}

// In PostgreSQL, only count() returns 0 for empty table.
// https://www.postgresql.org/docs/15/functions-aggregate.html
func (q *Queries) CountRowsEmptyTable(ctx context.Context) ([]CountRowsEmptyTableRow, error) {
	rows, err := q.db.Query(ctx, countRowsEmptyTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CountRowsEmptyTableRow
	for rows.Next() {
		var i CountRowsEmptyTableRow
		if err := rows.Scan(&i.A, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const firstRowFromEmptyTable = `-- name: FirstRowFromEmptyTable :many
SELECT a, (SELECT a FROM empty limit 1) as "first" FROM foo
`

type FirstRowFromEmptyTableRow struct {
	A     int32
	First pgtype.Int4
}

func (q *Queries) FirstRowFromEmptyTable(ctx context.Context) ([]FirstRowFromEmptyTableRow, error) {
	rows, err := q.db.Query(ctx, firstRowFromEmptyTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FirstRowFromEmptyTableRow
	for rows.Next() {
		var i FirstRowFromEmptyTableRow
		if err := rows.Scan(&i.A, &i.First); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const firstRowFromFooTable = `-- name: FirstRowFromFooTable :many
SELECT a, (SELECT a FROM foo limit 1) as "first" FROM foo
`

type FirstRowFromFooTableRow struct {
	A     int32
	First pgtype.Int4
}

func (q *Queries) FirstRowFromFooTable(ctx context.Context) ([]FirstRowFromFooTableRow, error) {
	rows, err := q.db.Query(ctx, firstRowFromFooTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FirstRowFromFooTableRow
	for rows.Next() {
		var i FirstRowFromFooTableRow
		if err := rows.Scan(&i.A, &i.First); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
