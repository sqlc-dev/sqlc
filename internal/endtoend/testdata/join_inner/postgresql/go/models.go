// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package querytest

import (
	"database/sql"
)

type Event struct {
	ID sql.NullInt32
}

type HandledEvent struct {
	LastHandledID sql.NullInt32
	Handler       sql.NullString
}
