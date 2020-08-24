package ast

import ()

type DeclareCursorStmt struct {
	Portalname *string
	Options    int
	Query      Node
}

func (n *DeclareCursorStmt) Pos() int {
	return 0
}
