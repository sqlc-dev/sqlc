// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package querytest

import (
	"database/sql"
)

type Bar struct {
	ID int32
}

type Foo struct {
	ID    int32
	BarID sql.NullInt32
}
