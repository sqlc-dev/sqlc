// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package querytest

import (
	"database/sql"
)

type Ft struct {
	B sql.NullString
}

type Tbl struct {
	A int64
	B sql.NullString
	C sql.NullString
	D sql.NullString
	E sql.NullInt64
}

type TblFt struct {
	B sql.NullString
	C sql.NullString
}
