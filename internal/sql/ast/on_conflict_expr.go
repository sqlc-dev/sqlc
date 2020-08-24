package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type OnConflictExpr struct {
	Action          OnConflictAction
	ArbiterElems    *List
	ArbiterWhere    ast.Node
	Constraint      Oid
	OnConflictSet   *List
	OnConflictWhere ast.Node
	ExclRelIndex    int
	ExclRelTlist    *List
}

func (n *OnConflictExpr) Pos() int {
	return 0
}
