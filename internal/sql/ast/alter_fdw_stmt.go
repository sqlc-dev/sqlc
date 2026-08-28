package ast

type AlterFdwStmt struct {
	Tag NodeTag[AlterFdwStmt] `json:"tag"`

	Fdwname     *string `json:"fdwname,omitempty"`
	FuncOptions *List   `json:"func_options,omitempty"`
	Options     *List   `json:"options,omitempty"`
}

func (n *AlterFdwStmt) Pos() int {
	return 0
}
