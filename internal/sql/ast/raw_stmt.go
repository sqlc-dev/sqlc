package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type RawStmt struct {
	Stmt         Node
	StmtLocation int
	StmtLen      int
}

func (n *RawStmt) Pos() int {
	return n.StmtLocation
}

func (n *RawStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n.Stmt != nil {
		buf.astFormat(n.Stmt, d)
	}
	buf.WriteString(";")
}
