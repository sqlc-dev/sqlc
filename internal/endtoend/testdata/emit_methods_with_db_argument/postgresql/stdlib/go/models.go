// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package querytest

import (
	"database/sql"
)

type User struct {
	ID        int32
	FirstName string
	LastName  sql.NullString
	Age       int32
}
