// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const aliasExpand = `-- name: AliasExpand :many
SELECT f.id, b.id, title
FROM foo f
JOIN bar b ON b.id = f.id
WHERE f.id = ?
`

type AliasExpandRow struct {
	ID    uint64
	ID_2  uint64
	Title sql.NullString
}

func (q *Queries) AliasExpand(ctx context.Context, id uint64, aq ...AdditionalQuery) ([]AliasExpandRow, error) {
	query := aliasExpand
	queryParams := []interface{}{id}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AliasExpandRow
	for rows.Next() {
		var i AliasExpandRow
		if err := rows.Scan(&i.ID, &i.ID_2, &i.Title); err != nil {
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

const aliasJoin = `-- name: AliasJoin :many
SELECT f.id, b.title
FROM foo f
JOIN bar b ON b.id = f.id
WHERE f.id = ?
`

type AliasJoinRow struct {
	ID    uint64
	Title sql.NullString
}

func (q *Queries) AliasJoin(ctx context.Context, id uint64, aq ...AdditionalQuery) ([]AliasJoinRow, error) {
	query := aliasJoin
	queryParams := []interface{}{id}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AliasJoinRow
	for rows.Next() {
		var i AliasJoinRow
		if err := rows.Scan(&i.ID, &i.Title); err != nil {
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
