package ast

type AlterFdwStmt struct {
	Fdwname     *string
	FuncOptions *List
	Options     *List
}

func (n *AlterFdwStmt) Pos() int {
	return 0
}
