package ast

type AlterTableMoveAllStmt struct {
	OrigTablespacename *string
	Objtype            ObjectType
	Roles              *List
	NewTablespacename  *string
	Nowait             bool
}

func (n *AlterTableMoveAllStmt) Pos() int {
	return 0
}
