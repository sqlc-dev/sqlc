package ast

type PrepareStmt struct {
	Tag NodeTag[PrepareStmt] `json:"tag"`

	Name     *string
	Argtypes *List
	Query    Node
}

func (n *PrepareStmt) Pos() int {
	return 0
}
