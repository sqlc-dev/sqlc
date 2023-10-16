// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getSomeDeletedNotOk = `-- name: GetSomeDeletedNotOk :many
DELETE FROM a
USING b
WHERE a.b_id_fk = b.b_id
RETURNING b.b_id
`

func (q *Queries) GetSomeDeletedNotOk(ctx context.Context) ([]pgtype.Text, error) {
	rows, err := q.db.Query(ctx, getSomeDeletedNotOk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.Text
	for rows.Next() {
		var b_id pgtype.Text
		if err := rows.Scan(&b_id); err != nil {
			return nil, err
		}
		items = append(items, b_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
