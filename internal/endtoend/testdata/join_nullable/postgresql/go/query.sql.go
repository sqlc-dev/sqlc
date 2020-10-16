// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const nullableJoin = `-- name: NullableJoin :many
SELECT f.id, f.bar_id, b.id
FROM foo f
FULL OUTER JOIN bar b ON b.id = f.bar_id
WHERE f.id = $1
`

type NullableJoinRow struct {
	ID    int32
	BarID sql.NullInt32
	ID_2  sql.NullInt32
}

func (q *Queries) NullableJoin(ctx context.Context, id int32) ([]NullableJoinRow, error) {
	rows, err := q.db.QueryContext(ctx, nullableJoin, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []NullableJoinRow
	for rows.Next() {
		var i NullableJoinRow
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
