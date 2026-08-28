package ast

type ExplainStmt struct {
	Tag NodeTag[ExplainStmt] `json:"tag"`

	Query   Node  `json:"query,omitempty"`
	Options *List `json:"options,omitempty"`
}

func (n *ExplainStmt) Pos() int {
	return 0
}
