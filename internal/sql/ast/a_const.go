package ast

type A_Const struct {
	Val      Node
	Location int
}

func (n *A_Const) Pos() int {
	return n.Location
}

func (n *A_Const) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if _, ok := n.Val.(*String); ok {
		buf.WriteString("'")
		buf.astFormat(n.Val)
		buf.WriteString("'")
	} else {
		buf.astFormat(n.Val)
	}
}
