package ast

type CompositeTypeStmt struct {
	TypeName *TypeName
	Cols     []*ColumnDef
}

func (n *CompositeTypeStmt) Pos() int {
	return 0
}
