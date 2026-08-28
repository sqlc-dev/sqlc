package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CommonTableExpr struct {
	Tag NodeTag[CommonTableExpr] `json:"tag"`

	Ctename          *string `json:"ctename,omitempty"`
	Aliascolnames    *List   `json:"aliascolnames,omitempty"`
	Ctequery         Node    `json:"ctequery,omitempty"`
	Location         int     `json:"location"`
	Cterecursive     bool    `json:"cterecursive"`
	Cterefcount      int     `json:"cterefcount"`
	Ctecolnames      *List   `json:"ctecolnames,omitempty"`
	Ctecoltypes      *List   `json:"ctecoltypes,omitempty"`
	Ctecoltypmods    *List   `json:"ctecoltypmods,omitempty"`
	Ctecolcollations *List   `json:"ctecolcollations,omitempty"`
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
