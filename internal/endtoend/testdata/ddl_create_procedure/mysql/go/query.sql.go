// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"
)

const callInsertData = `-- name: CallInsertData :exec
CALL insert_data(?, ?)
`

type CallInsertDataParams struct {
	A int32
	B int32
}

func (q *Queries) CallInsertData(ctx context.Context, arg CallInsertDataParams) error {
	_, err := q.db.ExecContext(ctx, callInsertData, arg.A, arg.B)
	return err
}

const callInsertDataNoArgs = `-- name: CallInsertDataNoArgs :exec
CALL insert_data(1, 2)
`

func (q *Queries) CallInsertDataNoArgs(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, callInsertDataNoArgs)
	return err
}

const callInsertDataSqlcArgs = `-- name: CallInsertDataSqlcArgs :exec
CALL insert_data(?, ?)
`

type CallInsertDataSqlcArgsParams struct {
	Foo int32
	Bar int32
}

func (q *Queries) CallInsertDataSqlcArgs(ctx context.Context, arg CallInsertDataSqlcArgsParams) error {
	_, err := q.db.ExecContext(ctx, callInsertDataSqlcArgs, arg.Foo, arg.Bar)
	return err
}
