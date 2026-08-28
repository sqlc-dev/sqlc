package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CaseWhen struct {
	Tag NodeTag[CaseWhen] `json:"tag"`

	Xpr      Node `json:"xpr,omitempty"`
	Expr     Node `json:"expr,omitempty"`
	Result   Node `json:"result,omitempty"`
	Location int  `json:"location"`
}

func (n *CaseWhen) Pos() int {
	return n.Location
}

func (n *CaseWhen) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("WHEN ")
	buf.astFormat(n.Expr, d)
	buf.WriteString(" THEN ")
	buf.astFormat(n.Result, d)
}
