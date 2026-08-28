package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type IndexElem struct {
	Tag NodeTag[IndexElem] `json:"tag"`

	Name          *string `json:",omitempty"`
	Expr          Node    `json:",omitempty"`
	Indexcolname  *string `json:",omitempty"`
	Collation     *List   `json:",omitempty"`
	Opclass       *List   `json:",omitempty"`
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
