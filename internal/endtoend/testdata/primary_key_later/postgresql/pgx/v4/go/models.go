// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package primary_key_later

import (
	"database/sql"
)

type Author struct {
	ID   int64
	Name string
	Bio  sql.NullString
}
