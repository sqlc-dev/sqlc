package ast

type A_Indices struct {
	IsSlice bool
	Lidx    Node
	Uidx    Node
}

func (n *A_Indices) Pos() int {
	return 0
}

func (n *A_Indices) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("[")
	if n.IsSlice {
		if set(n.Lidx) {
			buf.astFormat(n.Lidx)
		}
		buf.WriteString(":")
		if set(n.Uidx) {
			buf.astFormat(n.Uidx)
		}
	} else {
		buf.astFormat(n.Uidx)
	}
	buf.WriteString("]")
}
