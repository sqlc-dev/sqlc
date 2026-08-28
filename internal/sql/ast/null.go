package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type Null struct {
	Tag NodeTag[Null] `json:"tag"`
}

func (n *Null) Pos() int {
	return 0
}
func (n *Null) Format(buf *TrackedBuffer, d format.Dialect) {
	buf.WriteString("NULL")
}
