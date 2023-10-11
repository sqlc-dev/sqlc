package ast

type ResTarget struct {
	Name        *string
	Indirection *List
	Val         Node
	Location    int
}

func (n *ResTarget) Pos() int {
	return n.Location
}

func (n *ResTarget) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if set(n.Val) {
		buf.astFormat(n.Val)
		if n.Name != nil {
			buf.WriteString(" AS ")
			buf.WriteString(*n.Name)
		}
	} else {
		if n.Name != nil {
			buf.WriteString(*n.Name)
		}
	}
}
