// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package override

import (
	"context"

	t "github.com/jackc/pgtype"
)

const test = `-- name: test :exec
UPDATE foo SET langs = $1
`

func (q *Queries) test(ctx context.Context, langs *t.Text) error {
	query := test
	queryParams := []interface{}{langs}

	_, err := q.db.ExecContext(ctx, query, queryParams...)
	return err
}
