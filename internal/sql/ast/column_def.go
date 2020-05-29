package ast

type ColumnDef struct {
	Colname   string
	TypeName  *TypeName
	IsNotNull bool
	IsArray   bool
	Vals      *List
}

func (n *ColumnDef) Pos() int {
	return 0
}
