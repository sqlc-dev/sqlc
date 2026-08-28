package ast

type DeclareCursorStmt struct {
	Tag NodeTag[DeclareCursorStmt] `json:"tag"`

	Portalname *string
	Options    int
	Query      Node
}

func (n *DeclareCursorStmt) Pos() int {
	return 0
}
