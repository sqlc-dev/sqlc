package ast

type RawStmt struct {
	Stmt         Node
	StmtLocation int
	StmtLen      int
}

func (n *RawStmt) Pos() int {
	return n.StmtLocation
}
