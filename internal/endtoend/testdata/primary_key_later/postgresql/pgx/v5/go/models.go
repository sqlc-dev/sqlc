// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package primary_key_later

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Author struct {
	ID   int64
	Name string
	Bio  pgtype.Text
}
