package ast

type CreateFdwStmt struct {
	Tag NodeTag[CreateFdwStmt] `json:"tag"`

	Fdwname     *string
	FuncOptions *List
	Options     *List
}

func (n *CreateFdwStmt) Pos() int {
	return 0
}
