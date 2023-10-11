package ast

type A_Expr struct {
	Kind     A_Expr_Kind
	Name     *List
	Lexpr    Node
	Rexpr    Node
	Location int
}

func (n *A_Expr) Pos() int {
	return n.Location
}

func (n *A_Expr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.astFormat(n.Lexpr)
	buf.WriteString(" ")
	buf.astFormat(n.Name)
	buf.WriteString(" ")
	buf.astFormat(n.Rexpr)
}
