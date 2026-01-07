package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type String struct {
	Str string
}

func (n *String) Pos() int {
	return 0
}

func (n *String) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString(n.Str)
}
