package ast

type ColumnRef struct {
	Name string

	// From pg.ColumnRef
	Fields   *List
	Location int
}

func (n *ColumnRef) Pos() int {
	return n.Location
}

func (n *ColumnRef) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.astFormat(n.Fields)
}
