// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const listNullable = `-- name: ListNullable :many
SELECT
  NULL::text as a,
  NULL::integer as b,
  NULL::bigint as c,
  NULL::time as d
FROM foo
`

type ListNullableRow struct {
	A pgtype.Text
	B pgtype.Int4
	C pgtype.Int8
	D pgtype.Time
}

func (q *Queries) ListNullable(ctx context.Context) ([]ListNullableRow, error) {
	rows, err := q.db.Query(ctx, listNullable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListNullableRow
	for rows.Next() {
		var i ListNullableRow
		if err := rows.Scan(
			&i.A,
			&i.B,
			&i.C,
			&i.D,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
