package ast

type ColumnDef struct {
	Colname   string
	TypeName  *TypeName
	IsNotNull bool
	IsArray   bool
}

func (n *ColumnDef) Pos() int {
	return 0
}
