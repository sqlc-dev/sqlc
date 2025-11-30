package ast

type A_ArrayExpr struct {
	Elements *List
	Location int
}

func (n *A_ArrayExpr) Pos() int {
	return n.Location
}

func (n *A_ArrayExpr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("ARRAY[")
	buf.join(n.Elements, ", ")
	buf.WriteString("]")
}
