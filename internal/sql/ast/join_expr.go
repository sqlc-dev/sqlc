package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type JoinExpr struct {
	Jointype    JoinType
	IsNatural   bool
	Larg        ast.Node
	Rarg        ast.Node
	UsingClause *ast.List
	Quals       ast.Node
	Alias       *Alias
	Rtindex     int
}

func (n *JoinExpr) Pos() int {
	return 0
}
