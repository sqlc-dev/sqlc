// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package querytest

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Foo struct {
	Bar pgtype.Text
	Bat string
	Baz pgtype.Int8
	Qux int64
}
