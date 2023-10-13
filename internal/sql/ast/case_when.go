package ast

type CaseWhen struct {
	Xpr      Node
	Expr     Node
	Result   Node
	Location int
}

func (n *CaseWhen) Pos() int {
	return n.Location
}

func (n *CaseWhen) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("WHEN ")
	buf.astFormat(n.Expr)
	buf.WriteString(" THEN ")
	buf.astFormat(n.Result)
}
