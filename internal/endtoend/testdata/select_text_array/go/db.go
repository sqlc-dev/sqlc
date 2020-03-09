// Code generated by sqlc. DO NOT EDIT.

package querytest

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

type Querier interface {
	SelectTextArray(ctx context.Context, dollar_1 []string) ([][]string, error)
	WithTx(*sql.Tx) Querier
}

var _ Querier = (*Queries)(nil)

func (q *Queries) WithTx(tx *sql.Tx) Querier {
	return &Queries{
		db: tx,
	}
}
