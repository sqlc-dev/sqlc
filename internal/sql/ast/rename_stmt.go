package ast

type RenameStmt struct {
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
