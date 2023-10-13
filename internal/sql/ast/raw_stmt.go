package ast

type RawStmt struct {
	Stmt         Node
	StmtLocation int
	StmtLen      int
}

func (n *RawStmt) Pos() int {
	return n.StmtLocation
}

func (n *RawStmt) Format(buf *TrackedBuffer) {
	if n.Stmt != nil {
		buf.astFormat(n.Stmt)
	}
	buf.WriteString(";")
}
