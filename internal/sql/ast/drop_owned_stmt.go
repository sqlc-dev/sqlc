package ast

type DropOwnedStmt struct {
	Tag NodeTag[DropOwnedStmt] `json:"tag"`

	Roles    *List
	Behavior DropBehavior
}

func (n *DropOwnedStmt) Pos() int {
	return 0
}
