package ast

type InferClause struct {
	IndexElems  *List
	WhereClause Node
	Conname     *string
	Location    int
}

func (n *InferClause) Pos() int {
	return n.Location
}

func (n *InferClause) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if n.Conname != nil && *n.Conname != "" {
		buf.WriteString("ON CONSTRAINT ")
		buf.WriteString(*n.Conname)
	} else if items(n.IndexElems) {
		buf.WriteString("(")
		buf.join(n.IndexElems, ", ")
		buf.WriteString(")")
		if set(n.WhereClause) {
			buf.WriteString(" WHERE ")
			buf.astFormat(n.WhereClause)
		}
	}
}
