package ast

type PrepareStmt struct {
	Tag NodeTag[PrepareStmt] `json:"tag"`

	Name     *string `json:"name,omitempty"`
	Argtypes *List   `json:"argtypes,omitempty"`
	Query    Node    `json:"query,omitempty"`
}

func (n *PrepareStmt) Pos() int {
	return 0
}
