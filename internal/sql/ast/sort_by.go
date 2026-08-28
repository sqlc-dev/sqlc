package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type SortBy struct {
	Tag NodeTag[SortBy] `json:"tag"`

	Node        Node        `json:"node,omitempty"`
	SortbyDir   SortByDir   `json:"sortby_dir"`
	SortbyNulls SortByNulls `json:"sortby_nulls"`
	UseOp       *List       `json:"use_op,omitempty"`
	Location    int         `json:"location"`
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
