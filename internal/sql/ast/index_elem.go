package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type IndexElem struct {
	Tag NodeTag[IndexElem] `json:"tag"`

	Name          *string     `json:"name,omitempty"`
	Expr          Node        `json:"expr,omitempty"`
	Indexcolname  *string     `json:"indexcolname,omitempty"`
	Collation     *List       `json:"collation,omitempty"`
	Opclass       *List       `json:"opclass,omitempty"`
	Ordering      SortByDir   `json:"ordering"`
	NullsOrdering SortByNulls `json:"nulls_ordering"`
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
