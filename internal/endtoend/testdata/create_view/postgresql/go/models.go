// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package querytest

import (
	"database/sql"
)

type FirstView struct {
	Val string
}

type Foo struct {
	Val  string
	Val2 sql.NullInt32
}

type SecondView struct {
	Val  string
	Val2 sql.NullInt32
}
