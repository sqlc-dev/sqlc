// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package querytest

import (
	"database/sql"
)

type Foo struct {
	ID   int64
	Name string
	Bar  sql.NullString
}
