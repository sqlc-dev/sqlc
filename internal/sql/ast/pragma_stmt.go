package ast

// YDB specific
type Pragma_stmt struct {
	Name Node
	Cols *List
	Equals bool
	Values *List
	Location int
}

func (n *Pragma_stmt) Pos() int {
	return n.Location
}

func (n *Pragma_stmt) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}

	buf.WriteString("PRAGMA ")
	if n.Name != nil {
		buf.astFormat(n.Name)
	}
	if n.Cols != nil {
		buf.astFormat(n.Cols)
	}

	if n.Equals {
		buf.WriteString(" = ")
	}

	if n.Values != nil {
		if n.Equals {
			buf.astFormat(n.Values)
		} else {
			buf.WriteString("(")
			buf.astFormat(n.Values)
			buf.WriteString(")")
		}
	}

}
