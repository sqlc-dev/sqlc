package ast

type CoalesceExpr struct {
	Xpr            Node
	Coalescetype   Oid
	Coalescecollid Oid
	Args           *List
	Location       int
}

func (n *CoalesceExpr) Pos() int {
	return n.Location
}

func (n *CoalesceExpr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("COALESCE(")
	buf.astFormat(n.Args)
	buf.WriteString(")")
}
