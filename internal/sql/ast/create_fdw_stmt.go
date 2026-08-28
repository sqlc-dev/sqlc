package ast

type CreateFdwStmt struct {
	Tag NodeTag[CreateFdwStmt] `json:"tag"`

	Fdwname     *string `json:",omitempty"`
	FuncOptions *List   `json:",omitempty"`
	Options     *List   `json:",omitempty"`
}

func (n *CreateFdwStmt) Pos() int {
	return 0
}
