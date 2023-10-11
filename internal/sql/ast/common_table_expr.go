package ast

import (
	"fmt"
)

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
		fmt.Fprintf(buf, " %s AS (", *n.Ctename)
	}
	buf.astFormat(n.Ctequery)
	buf.WriteString(")")
}
