package ast

import ()

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
