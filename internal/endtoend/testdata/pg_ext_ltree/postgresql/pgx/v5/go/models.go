// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package querytest

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Foo struct {
	QualifiedName pgtype.Text
	NameQuery     pgtype.Text
	FtsNameQuery  pgtype.Text
}
