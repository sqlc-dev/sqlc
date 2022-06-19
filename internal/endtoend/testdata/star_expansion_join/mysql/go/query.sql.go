// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const starExpansionJoin = `-- name: StarExpansionJoin :many
SELECT a, b, c, d FROM foo, bar
`

type StarExpansionJoinRow struct {
	A sql.NullString
	B sql.NullString
	C sql.NullString
	D sql.NullString
}

func (q *Queries) StarExpansionJoin(ctx context.Context) ([]StarExpansionJoinRow, error) {
	ctx, done := q.observer(ctx, "StarExpansionJoin")
	rows, err := q.db.QueryContext(ctx, starExpansionJoin)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []StarExpansionJoinRow
	for rows.Next() {
		var i StarExpansionJoinRow
		if err := rows.Scan(
			&i.A,
			&i.B,
			&i.C,
			&i.D,
		); err != nil {
			return nil, done(err)
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}
