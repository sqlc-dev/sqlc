package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CommonTableExpr struct {
	Tag NodeTag[CommonTableExpr] `json:"tag"`

	Ctename          *string `json:",omitempty"`
	Aliascolnames    *List   `json:",omitempty"`
	Ctequery         Node    `json:",omitempty"`
	Location         int
	Cterecursive     bool
	Cterefcount      int
	Ctecolnames      *List `json:",omitempty"`
	Ctecoltypes      *List `json:",omitempty"`
	Ctecoltypmods    *List `json:",omitempty"`
	Ctecolcollations *List `json:",omitempty"`
}

func (n *CommonTableExpr) Pos() int {
	return n.Location
}

func (n *CommonTableExpr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if n.Ctename != nil {
		buf.WriteString(*n.Ctename)
	}
	if items(n.Aliascolnames) {
		buf.WriteString("(")
		buf.join(n.Aliascolnames, d, ", ")
		buf.WriteString(")")
	}
	buf.WriteString(" AS (")
	buf.Group()
	buf.Indent()
	// The body's boundary: a comment at the top of the CTE prints here, and
	// an author who broke the line after the parenthesis keeps the CTE — and
	// therefore the statement around it — broken open.
	buf.boundary(n.Ctequery)
	buf.Softline()
	buf.astFormat(n.Ctequery, d)
	buf.EndIndent()
	buf.Softline()
	buf.EndGroup()
	buf.WriteString(")")
}
