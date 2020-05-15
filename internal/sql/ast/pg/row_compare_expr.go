package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RowCompareExpr struct {
	Xpr          ast.Node
	Rctype       RowCompareType
	Opnos        *ast.List
	Opfamilies   *ast.List
	Inputcollids *ast.List
	Largs        *ast.List
	Rargs        *ast.List
}

func (n *RowCompareExpr) Pos() int {
	return 0
}
