package ast

type RawStmt_PG struct {
	Stmt         Node
	StmtLocation int
	StmtLen      int
}

func (n *RawStmt_PG) Pos() int {
	return 0
}
