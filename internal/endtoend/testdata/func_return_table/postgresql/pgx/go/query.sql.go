// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const foo = `-- name: Foo :one
SELECT register_account FROM register_account('a', 'b')
`

func (q *Queries) Foo(ctx context.Context) (pgtype.Int4, error) {
	row := q.db.QueryRow(ctx, foo)
	var register_account pgtype.Int4
	err := row.Scan(&register_account)
	return register_account, err
}
