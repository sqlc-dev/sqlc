// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package querytest

import (
	"context"
)

type Querier interface {
	InsertBar(ctx context.Context) (int64, error)
}

var _ Querier = (*Queries)(nil)
