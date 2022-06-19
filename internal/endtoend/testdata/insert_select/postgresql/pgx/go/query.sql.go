// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
)

const insertSelect = `-- name: InsertSelect :exec
INSERT INTO foo (name, meta)
SELECT name, $1
FROM bar WHERE ready = $2
`

type InsertSelectParams struct {
	Meta  string
	Ready bool
}

func (q *Queries) InsertSelect(ctx context.Context, arg InsertSelectParams) error {
	ctx, done := q.observer(ctx, "InsertSelect")
	_, err := q.db.Exec(ctx, insertSelect, arg.Meta, arg.Ready)
	return done(err)
}
