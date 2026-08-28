package ast

type AlterFdwStmt struct {
	Tag NodeTag[AlterFdwStmt] `json:"tag"`

	Fdwname     *string `json:",omitempty"`
	FuncOptions *List   `json:",omitempty"`
	Options     *List   `json:",omitempty"`
}

func (n *AlterFdwStmt) Pos() int {
	return 0
}
