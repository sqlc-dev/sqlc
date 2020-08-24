package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CommonTableExpr struct {
	Ctename          *string
	Aliascolnames    *List
	Ctequery         ast.Node
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
