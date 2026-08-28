package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type RowExpr struct {
	Tag NodeTag[RowExpr] `json:"tag"`

	Xpr       Node         `json:"xpr,omitempty"`
	Args      *List        `json:"args,omitempty"`
	RowTypeid Oid          `json:"row_typeid"`
	RowFormat CoercionForm `json:"row_format"`
	Colnames  *List        `json:"colnames,omitempty"`
	Location  int          `json:"location"`
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
