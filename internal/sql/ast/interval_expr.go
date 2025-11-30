package ast

// IntervalExpr represents a MySQL INTERVAL expression like "INTERVAL 1 DAY"
type IntervalExpr struct {
	Value    Node
	Unit     string
	Location int
}

func (n *IntervalExpr) Pos() int {
	return n.Location
}

func (n *IntervalExpr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("INTERVAL ")
	buf.astFormat(n.Value)
	buf.WriteString(" ")
	buf.WriteString(n.Unit)
}
