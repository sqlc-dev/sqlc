package ast

type DiscardStmt struct {
	Tag NodeTag[DiscardStmt] `json:"tag"`

	Target DiscardMode `json:"target"`
}

func (n *DiscardStmt) Pos() int {
	return 0
}
