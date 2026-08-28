package ast

type DropStmt struct {
	Tag NodeTag[DropStmt] `json:"tag"`

	Objects    *List
	RemoveType ObjectType
	Behavior   DropBehavior
	MissingOk  bool
	Concurrent bool
}

func (n *DropStmt) Pos() int {
	return 0
}
