// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package querytest

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Activity struct {
	AccountID int64
	EventTime pgtype.Timestamptz
}
