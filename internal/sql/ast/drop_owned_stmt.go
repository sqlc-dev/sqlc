package ast

type DropOwnedStmt struct {
	Tag NodeTag[DropOwnedStmt] `json:"tag"`

	Roles    *List `json:",omitempty"`
	Behavior DropBehavior
}

func (n *DropOwnedStmt) Pos() int {
	return 0
}
