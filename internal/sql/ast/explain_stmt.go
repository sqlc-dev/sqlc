package ast

type ExplainStmt struct {
	Tag NodeTag[ExplainStmt] `json:"tag"`

	Query   Node  `json:",omitempty"`
	Options *List `json:",omitempty"`
}

func (n *ExplainStmt) Pos() int {
	return 0
}
