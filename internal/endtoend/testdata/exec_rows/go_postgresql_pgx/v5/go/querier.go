// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Querier interface {
	WithTx(tx pgx.Tx) *Queries
	DeleteBarByID(ctx context.Context, id int32) (int64, error)
}

var _ Querier = (*Queries)(nil)
