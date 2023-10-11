package ast

type List struct {
	Items []Node
}

func (n *List) Pos() int {
	return 0
}

func (n *List) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.join(n, ",")
}
