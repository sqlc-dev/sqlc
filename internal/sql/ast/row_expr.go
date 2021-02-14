package ast

type RowExpr struct {
	Xpr       Node
	Args      *List
	RowTypeid Oid
	RowFormat CoercionForm
	Colnames  *List
	Location  int
}

func (n *RowExpr) Pos() int {
	return n.Location
}
