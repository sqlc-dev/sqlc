package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

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

func (n *RowExpr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if items(n.Args) {
		buf.WriteString("args")
		buf.astFormat(n.Args, d)
	}
	buf.astFormat(n.Xpr, d)
	if items(n.Colnames) {
		buf.WriteString("cols")
		buf.astFormat(n.Colnames, d)
	}
}
