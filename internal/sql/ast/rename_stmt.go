package ast

type RenameStmt struct {
	Tag NodeTag[RenameStmt] `json:"tag"`

	RenameType   ObjectType
	RelationType ObjectType
	Relation     *RangeVar `json:",omitempty"`
	Object       Node      `json:",omitempty"`
	Subname      *string   `json:",omitempty"`
	Newname      *string   `json:",omitempty"`
	Behavior     DropBehavior
	MissingOk    bool
}

func (n *RenameStmt) Pos() int {
	return 0
}
