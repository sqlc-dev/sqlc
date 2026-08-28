package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CaseExpr struct {
	Tag NodeTag[CaseExpr] `json:"tag"`

	Xpr        Node  `json:"xpr,omitempty"`
	Casetype   Oid   `json:"casetype"`
	Casecollid Oid   `json:"casecollid"`
	Arg        Node  `json:"arg,omitempty"`
	Args       *List `json:"args,omitempty"`
	Defresult  Node  `json:"defresult,omitempty"`
	Location   int   `json:"location"`
}

func (n *CaseExpr) Pos() int {
	return n.Location
}

func (n *CaseExpr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("CASE ")
	if set(n.Arg) {
		buf.astFormat(n.Arg, d)
		buf.WriteString(" ")
	}
	buf.join(n.Args, d, " ")
	if set(n.Defresult) {
		buf.WriteString(" ELSE ")
		buf.astFormat(n.Defresult, d)
	}
	buf.WriteString(" END")
}
