// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const starExpansionCTE = `-- name: StarExpansionCTE :many
WITH cte AS (SELECT a, b FROM foo) SELECT c, d FROM bar
`

func (q *Queries) StarExpansionCTE(ctx context.Context) ([]Bar, error) {
	rows, err := q.db.Query(ctx, starExpansionCTE)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Bar
	for rows.Next() {
		var i Bar
		if err := rows.Scan(&i.C, &i.D); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const starExpansionTwoCTE = `-- name: StarExpansionTwoCTE :many
WITH 
  a AS (SELECT a, b FROM foo),
  b AS (SELECT 1::int as bar, a, b FROM a)
SELECT bar, a, b FROM b
`

type StarExpansionTwoCTERow struct {
	Bar int32
	A   sql.NullString
	B   sql.NullString
}

func (q *Queries) StarExpansionTwoCTE(ctx context.Context) ([]StarExpansionTwoCTERow, error) {
	rows, err := q.db.Query(ctx, starExpansionTwoCTE)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []StarExpansionTwoCTERow
	for rows.Next() {
		var i StarExpansionTwoCTERow
		if err := rows.Scan(&i.Bar, &i.A, &i.B); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
