package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CaseExpr struct {
	Xpr        Node
	Casetype   Oid
	Casecollid Oid
	Arg        Node
	Args       *List
	Defresult  Node
	Location   int
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
