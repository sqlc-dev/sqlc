package ast

// TODO: Support array types
type ColumnDef struct {
	Colname   string
	TypeName  *TypeName
	IsNotNull bool
}

func (n *ColumnDef) Pos() int {
	return 0
}
