package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CommonTableExpr struct {
	Ctename          *string
	Aliascolnames    *ast.List
	Ctequery         ast.Node
	Location         int
	Cterecursive     bool
	Cterefcount      int
	Ctecolnames      *ast.List
	Ctecoltypes      *ast.List
	Ctecoltypmods    *ast.List
	Ctecolcollations *ast.List
}

func (n *CommonTableExpr) Pos() int {
	return n.Location
}
