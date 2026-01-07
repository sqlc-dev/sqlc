package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type InferClause struct {
	IndexElems  *List
	WhereClause Node
	Conname     *string
	Location    int
}

func (n *InferClause) Pos() int {
	return n.Location
}

func (n *InferClause) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if n.Conname != nil && *n.Conname != "" {
		buf.WriteString("ON CONSTRAINT ")
		buf.WriteString(*n.Conname)
	} else if items(n.IndexElems) {
		buf.WriteString("(")
		buf.join(n.IndexElems, d, ", ")
		buf.WriteString(")")
		if set(n.WhereClause) {
			buf.WriteString(" WHERE ")
			buf.astFormat(n.WhereClause, d)
		}
	}
}
