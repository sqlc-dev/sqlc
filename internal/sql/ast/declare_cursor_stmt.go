package ast

type DeclareCursorStmt struct {
	Tag NodeTag[DeclareCursorStmt] `json:"tag"`

	Portalname *string `json:",omitempty"`
	Options    int
	Query      Node `json:",omitempty"`
}

func (n *DeclareCursorStmt) Pos() int {
	return 0
}
