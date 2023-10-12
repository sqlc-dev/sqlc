package ast

type RangeSubselect struct {
	Lateral  bool
	Subquery Node
	Alias    *Alias
}

func (n *RangeSubselect) Pos() int {
	return 0
}

func (n *RangeSubselect) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("(")
	buf.astFormat(n.Subquery)
	buf.WriteString(")")
	if n.Alias != nil {
		buf.WriteString(" ")
		buf.astFormat(n.Alias)
	}
}
