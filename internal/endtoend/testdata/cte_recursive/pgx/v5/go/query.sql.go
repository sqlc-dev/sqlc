// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const cTERecursive = `-- name: CTERecursive :many
WITH RECURSIVE cte AS (
        SELECT b.id, b.parent_id FROM bar AS b
        WHERE b.id = $1
    UNION ALL
        SELECT b.id, b.parent_id
        FROM bar AS b, cte AS c
        WHERE b.parent_id = c.id
) SELECT id, parent_id FROM cte
`

type CTERecursiveRow struct {
	ID       int32
	ParentID pgtype.Int4
}

func (q *Queries) CTERecursive(ctx context.Context, id int32, aq ...AdditionalQuery) ([]CTERecursiveRow, error) {
	query := cTERecursive
	queryParams := []interface{}{id}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	rows, err := q.db.Query(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CTERecursiveRow
	for rows.Next() {
		var i CTERecursiveRow
		if err := rows.Scan(&i.ID, &i.ParentID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
