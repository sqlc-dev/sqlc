// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package querytest

import (
	"database/sql"
)

type Customer struct {
	CustID   int64
	CustName sql.NullString
	CustAddr sql.NullString
}

type CustomerAddress struct {
	CustID   int64
	CustAddr sql.NullString
}
