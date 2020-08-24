package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type UpdateStmt struct {
	Relation      *RangeVar
	TargetList    *ast.List
	WhereClause   ast.Node
	FromClause    *ast.List
	ReturningList *ast.List
	WithClause    *WithClause
}

func (n *UpdateStmt) Pos() int {
	return 0
}
