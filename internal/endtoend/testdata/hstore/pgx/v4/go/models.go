// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.1

package hstore

import (
	"github.com/jackc/pgtype"
)

type Foo struct {
	Bar pgtype.Hstore
	Baz pgtype.Hstore
}
