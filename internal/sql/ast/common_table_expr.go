package ast

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

func (n *CommonTableExpr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if n.Ctename != nil {
		buf.WriteString(*n.Ctename)
	}
	if items(n.Aliascolnames) {
		buf.WriteString("(")
		buf.join(n.Aliascolnames, ", ")
		buf.WriteString(")")
	}
	buf.WriteString(" AS (")
	buf.astFormat(n.Ctequery)
	buf.WriteString(")")
}
