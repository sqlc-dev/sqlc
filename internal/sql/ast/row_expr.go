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

func (n *RowExpr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if items(n.Args) {
		buf.WriteString("args")
		buf.astFormat(n.Args)
	}
	buf.astFormat(n.Xpr)
	if items(n.Colnames) {
		buf.WriteString("cols")
		buf.astFormat(n.Colnames)
	}
}
