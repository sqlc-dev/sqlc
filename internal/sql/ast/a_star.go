package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type A_Star struct {
}

func (n *A_Star) Pos() int {
	return 0
}

func (n *A_Star) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteRune('*')
}
