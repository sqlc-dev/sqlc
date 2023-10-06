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
	buf.astFormat(n.Val)
}
