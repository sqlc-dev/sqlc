package ast

type AlterFdwStmt struct {
	Tag NodeTag[AlterFdwStmt] `json:"tag"`

	Fdwname     *string
	FuncOptions *List
	Options     *List
}

func (n *AlterFdwStmt) Pos() int {
	return 0
}
