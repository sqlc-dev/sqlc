// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package querytest

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Bar struct {
	ColA pgtype.Bool
	ColB pgtype.Bool
}

type Foo struct {
	ColA bool
	ColB bool
}
