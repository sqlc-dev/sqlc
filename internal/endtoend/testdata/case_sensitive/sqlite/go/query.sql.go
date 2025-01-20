// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const insertContact = `-- name: InsertContact :exec
INSERT INTO contacts (
    pid,
    CustomerName
)
VALUES (?,?)
`

type InsertContactParams struct {
	Pid          sql.NullString
	Customername sql.NullString
}

func (q *Queries) InsertContact(ctx context.Context, arg InsertContactParams) error {
	_, err := q.db.ExecContext(ctx, insertContact, arg.Pid, arg.Customername)
	return err
}
