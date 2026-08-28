package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CaseWhen struct {
	Tag NodeTag[CaseWhen] `json:"tag"`

	Xpr      Node `json:",omitempty"`
	Expr     Node `json:",omitempty"`
	Result   Node `json:",omitempty"`
	Location int
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
