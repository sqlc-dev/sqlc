package ast

type RenameStmt struct {
	Tag NodeTag[RenameStmt] `json:"tag"`

	RenameType   ObjectType   `json:"rename_type"`
	RelationType ObjectType   `json:"relation_type"`
	Relation     *RangeVar    `json:"relation,omitempty"`
	Object       Node         `json:"object,omitempty"`
	Subname      *string      `json:"subname,omitempty"`
	Newname      *string      `json:"newname,omitempty"`
	Behavior     DropBehavior `json:"behavior"`
	MissingOk    bool         `json:"missing_ok"`
}

func (n *RenameStmt) Pos() int {
	return 0
}
