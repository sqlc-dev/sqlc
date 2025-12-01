package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type MultiAssignRef struct {
	Source   Node
	Colno    int
	Ncolumns int
}

func (n *MultiAssignRef) Pos() int {
	return 0
}

func (n *MultiAssignRef) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.astFormat(n.Source, d)
}
