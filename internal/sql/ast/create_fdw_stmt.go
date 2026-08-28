package ast

type CreateFdwStmt struct {
	Tag NodeTag[CreateFdwStmt] `json:"tag"`

	Fdwname     *string `json:"fdwname,omitempty"`
	FuncOptions *List   `json:"func_options,omitempty"`
	Options     *List   `json:"options,omitempty"`
}

func (n *CreateFdwStmt) Pos() int {
	return 0
}
