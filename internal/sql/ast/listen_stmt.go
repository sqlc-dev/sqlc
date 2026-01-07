package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type ListenStmt struct {
	Conditionname *string
}

func (n *ListenStmt) Pos() int {
	return 0
}

func (n *ListenStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("LISTEN ")
	if n.Conditionname != nil {
		buf.WriteString(*n.Conditionname)
	}
}
