// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.1

package authors

import (
	"database/sql"
)

type Author struct {
	ID   int64
	Name string
	Bio  sql.NullString
}

type Book struct {
	ID    int64
	Title string
}
