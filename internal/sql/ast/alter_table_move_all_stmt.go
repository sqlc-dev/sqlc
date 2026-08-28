package ast

type AlterTableMoveAllStmt struct {
	Tag NodeTag[AlterTableMoveAllStmt] `json:"tag"`

	OrigTablespacename *string    `json:"orig_tablespacename,omitempty"`
	Objtype            ObjectType `json:"objtype"`
	Roles              *List      `json:"roles,omitempty"`
	NewTablespacename  *string    `json:"new_tablespacename,omitempty"`
	Nowait             bool       `json:"nowait"`
}

func (n *AlterTableMoveAllStmt) Pos() int {
	return 0
}
