// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package querytest

import (
	"context"
	"database/sql"
)

type Querier interface {
	GetAll(ctx context.Context, aq ...AdditionalQuery) ([]*Foo, error)
	GetAllAByB(ctx context.Context, b sql.NullInt32, aq ...AdditionalQuery) ([]sql.NullInt32, error)
	GetOne(ctx context.Context, arg *GetOneParams, aq ...AdditionalQuery) (*Foo, error)
}

var _ Querier = (*Queries)(nil)
