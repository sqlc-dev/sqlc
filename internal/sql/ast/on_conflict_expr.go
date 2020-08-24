package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type OnConflictExpr struct {
	Action          OnConflictAction
	ArbiterElems    *ast.List
	ArbiterWhere    ast.Node
	Constraint      Oid
	OnConflictSet   *ast.List
	OnConflictWhere ast.Node
	ExclRelIndex    int
	ExclRelTlist    *ast.List
}

func (n *OnConflictExpr) Pos() int {
	return 0
}
