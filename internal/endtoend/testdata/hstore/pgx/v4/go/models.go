// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package hstore

import (
	"github.com/jackc/pgtype"
)

type Foo struct {
	Bar pgtype.Hstore
	Baz pgtype.Hstore
}
