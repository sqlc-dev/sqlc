// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package querytest

import (
	"github.com/jackc/pgtype"
)

type Transaction struct {
	ID        int64
	Uri       string
	ProgramID string
	Data      pgtype.JSONB
}
