package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type IndexElem struct {
	Name          *string
	Expr          Node
	Indexcolname  *string
	Collation     *List
	Opclass       *List
	Ordering      SortByDir
	NullsOrdering SortByNulls
}

func (n *IndexElem) Pos() int {
	return 0
}

func (n *IndexElem) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if n.Name != nil && *n.Name != "" {
		buf.WriteString(*n.Name)
	} else if set(n.Expr) {
		buf.astFormat(n.Expr, d)
	}
}
