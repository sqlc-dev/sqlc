package ast

type ExplainStmt struct {
	Tag NodeTag[ExplainStmt] `json:"tag"`

	Query   Node
	Options *List
}

func (n *ExplainStmt) Pos() int {
	return 0
}
