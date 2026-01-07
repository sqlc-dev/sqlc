package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CommonTableExpr struct {
	Ctename          *string
	Aliascolnames    *List
	Ctequery         Node
	Location         int
	Cterecursive     bool
	Cterefcount      int
	Ctecolnames      *List
	Ctecoltypes      *List
	Ctecoltypmods    *List
	Ctecolcollations *List
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
	buf.astFormat(n.Ctequery, d)
	buf.WriteString(")")
}
