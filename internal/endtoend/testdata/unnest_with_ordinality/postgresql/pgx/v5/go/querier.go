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
	GetValues(ctx context.Context) ([]GetValuesRow, error)
}

var _ Querier = (*Queries)(nil)
