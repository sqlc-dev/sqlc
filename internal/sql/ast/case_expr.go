package ast

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

func (n *CaseExpr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("CASE ")
	if set(n.Arg) {
		buf.astFormat(n.Arg)
		buf.WriteString(" ")
	}
	buf.join(n.Args, " ")
	if set(n.Defresult) {
		buf.WriteString(" ELSE ")
		buf.astFormat(n.Defresult)
	}
	buf.WriteString(" END")
}
