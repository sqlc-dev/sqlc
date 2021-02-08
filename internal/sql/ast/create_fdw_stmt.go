package ast

type CreateFdwStmt struct {
	Fdwname     *string
	FuncOptions *List
	Options     *List
}

func (n *CreateFdwStmt) Pos() int {
	return 0
}
