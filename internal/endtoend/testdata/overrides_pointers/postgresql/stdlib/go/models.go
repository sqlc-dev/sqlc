// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package override

import (
	t "github.com/jackc/pgtype"
)

type Foo struct {
	Other   string
	Total   int64
	Tags    []string
	ByteSeq []byte
	Retyped string
	Langs   *t.Text
}
