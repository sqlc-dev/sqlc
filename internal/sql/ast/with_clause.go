package ast

type WithClause struct {
	Ctes      *List
	Recursive bool
	Location  int
}

func (n *WithClause) Pos() int {
	return n.Location
}

func (n *WithClause) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("WITH")
	buf.astFormat(n.Ctes)
}
