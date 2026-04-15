package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type Float struct {
	Str string
}

func (n *Float) Pos() int {
	return 0
}

func (n *Float) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString(n.Str)
}
