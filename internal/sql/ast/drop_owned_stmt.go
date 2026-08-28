package ast

type DropOwnedStmt struct {
	Tag NodeTag[DropOwnedStmt] `json:"tag"`

	Roles    *List        `json:"roles,omitempty"`
	Behavior DropBehavior `json:"behavior"`
}

func (n *DropOwnedStmt) Pos() int {
	return 0
}
