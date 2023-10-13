package ast

type NamedArgExpr struct {
	Xpr       Node
	Arg       Node
	Name      *string
	Argnumber int
	Location  int
}

func (n *NamedArgExpr) Pos() int {
	return n.Location
}

func (n *NamedArgExpr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if n.Name != nil {
		buf.WriteString(*n.Name)
	}
	buf.WriteString(" => ")
	buf.astFormat(n.Arg)
}
