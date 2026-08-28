package ast

type DropStmt struct {
	Tag NodeTag[DropStmt] `json:"tag"`

	Objects    *List        `json:"objects,omitempty"`
	RemoveType ObjectType   `json:"remove_type"`
	Behavior   DropBehavior `json:"behavior"`
	MissingOk  bool         `json:"missing_ok"`
	Concurrent bool         `json:"concurrent"`
}

func (n *DropStmt) Pos() int {
	return 0
}
