// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package querytest

import (
	"database/sql"
)

type User struct {
	FirstName sql.NullString `json:"FirstName"`
	LastName  sql.NullString `json:"LastName"`
	Age       sql.NullInt16  `json:"Age"`
}
