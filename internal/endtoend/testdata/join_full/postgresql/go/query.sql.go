// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const fullJoin = `-- name: FullJoin :many
SELECT f.id, f.bar_id, b.id
FROM foo f
FULL OUTER JOIN bar b ON b.id = f.bar_id
WHERE f.id = $1
`

type FullJoinRow struct {
	ID    sql.NullInt32
	BarID sql.NullInt32
	ID_2  sql.NullInt32
}

func (q *Queries) FullJoin(ctx context.Context, id int32, aq ...AdditionalQuery) ([]FullJoinRow, error) {
	query := fullJoin
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
	var items []FullJoinRow
	for rows.Next() {
		var i FullJoinRow
		if err := rows.Scan(&i.ID, &i.BarID, &i.ID_2); err != nil {
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
