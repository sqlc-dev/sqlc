// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package querytest

import (
	"context"
	"database/sql"
)

type Querier interface {
	GetAll(ctx context.Context) ([]*Foo, error)
	GetAllAByB(ctx context.Context, b sql.NullInt32) ([]sql.NullInt32, error)
	GetOne(ctx context.Context, arg *GetOneParams) (*Foo, error)
}

var _ Querier = (*Queries)(nil)
