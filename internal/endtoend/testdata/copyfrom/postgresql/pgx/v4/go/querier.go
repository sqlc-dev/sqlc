// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package querytest

import (
	"context"
	"database/sql"
)

type Querier interface {
	// InsertSingleValue inserts a single value using copy.
	InsertSingleValue(ctx context.Context, a []sql.NullString) (int64, error)
	// InsertValues inserts multiple values using copy.
	InsertValues(ctx context.Context, arg []InsertValuesParams) (int64, error)
}

var _ Querier = (*Queries)(nil)
