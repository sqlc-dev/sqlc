package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type RawStmt struct {
	Tag NodeTag[RawStmt] `json:"tag"`

	Stmt         Node `json:",omitempty"`
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
	// The terminator goes first: a trailing line comment would swallow it.
	buf.WriteString(";")
	buf.flushRemaining()
}
