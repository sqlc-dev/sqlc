package ast

type RenameStmt struct {
	Tag NodeTag[RenameStmt] `json:"tag"`

	RenameType   ObjectType
	RelationType ObjectType
	Relation     *RangeVar
	Object       Node
	Subname      *string
	Newname      *string
	Behavior     DropBehavior
	MissingOk    bool
}

func (n *RenameStmt) Pos() int {
	return 0
}
