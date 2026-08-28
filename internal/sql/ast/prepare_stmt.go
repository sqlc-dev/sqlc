package ast

type PrepareStmt struct {
	Tag NodeTag[PrepareStmt] `json:"tag"`

	Name     *string `json:",omitempty"`
	Argtypes *List   `json:",omitempty"`
	Query    Node    `json:",omitempty"`
}

func (n *PrepareStmt) Pos() int {
	return 0
}
