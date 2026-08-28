package ast

type DeclareCursorStmt struct {
	Tag NodeTag[DeclareCursorStmt] `json:"tag"`

	Portalname *string `json:"portalname,omitempty"`
	Options    int     `json:"options"`
	Query      Node    `json:"query,omitempty"`
}

func (n *DeclareCursorStmt) Pos() int {
	return 0
}
