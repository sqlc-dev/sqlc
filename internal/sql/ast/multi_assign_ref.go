package ast

type MultiAssignRef struct {
	Source   Node
	Colno    int
	Ncolumns int
}

func (n *MultiAssignRef) Pos() int {
	return 0
}

func (n *MultiAssignRef) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.astFormat(n.Source)
}
