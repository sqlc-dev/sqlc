// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package querytest

import (
	"database/sql"
)

type User struct {
	ID        int64
	FirstName string
	LastName  sql.NullString
	Age       int64
}
