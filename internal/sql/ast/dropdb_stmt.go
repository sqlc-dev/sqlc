package ast

type DropdbStmt struct {
	Dbname    *string
	MissingOk bool
}

func (n *DropdbStmt) Pos() int {
	return 0
}
