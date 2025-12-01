package ast

type CollateExpr struct {
	Xpr      Node
	Arg      Node
	CollOid  Oid
	Location int
}

func (n *CollateExpr) Pos() int {
	return n.Location
}

func (n *CollateExpr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.astFormat(n.Xpr)
	buf.WriteString(" COLLATE ")
	buf.astFormat(n.Arg)
}
