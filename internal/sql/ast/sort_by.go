package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type SortBy struct {
	Node        Node
	SortbyDir   SortByDir
	SortbyNulls SortByNulls
	UseOp       *List
	Location    int
}

func (n *SortBy) Pos() int {
	return n.Location
}

func (n *SortBy) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.astFormat(n.Node, d)
	switch n.SortbyDir {
	case SortByDirAsc:
		buf.WriteString(" ASC")
	case SortByDirDesc:
		buf.WriteString(" DESC")
	}
	switch n.SortbyNulls {
	case SortByNullsFirst:
		buf.WriteString(" NULLS FIRST")
	case SortByNullsLast:
		buf.WriteString(" NULLS LAST")
	}
}
