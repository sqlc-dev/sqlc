package ast

type CompositeTypeStmt struct {
	Tag NodeTag[CompositeTypeStmt] `json:"tag"`

	TypeName *TypeName `json:"type_name,omitempty"`
	// Coldeflist is the type's fields, each a ColumnDef.
	Coldeflist *List `json:"coldeflist,omitempty"`
}

func (n *CompositeTypeStmt) Pos() int {
	return 0
}
