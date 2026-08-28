package ast

type AlterTableMoveAllStmt struct {
	Tag NodeTag[AlterTableMoveAllStmt] `json:"tag"`

	OrigTablespacename *string `json:",omitempty"`
	Objtype            ObjectType
	Roles              *List   `json:",omitempty"`
	NewTablespacename  *string `json:",omitempty"`
	Nowait             bool
}

func (n *AlterTableMoveAllStmt) Pos() int {
	return 0
}
