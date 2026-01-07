package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CreateExtensionStmt struct {
	Extname     *string
	IfNotExists bool
	Options     *List
}

func (n *CreateExtensionStmt) Pos() int {
	return 0
}

func (n *CreateExtensionStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("CREATE EXTENSION ")
	if n.IfNotExists {
		buf.WriteString("IF NOT EXISTS ")
	}
	if n.Extname != nil {
		buf.WriteString(*n.Extname)
	}
}
