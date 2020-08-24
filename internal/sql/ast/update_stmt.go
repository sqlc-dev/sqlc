package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type UpdateStmt struct {
	Relation      *RangeVar
	TargetList    *List
	WhereClause   ast.Node
	FromClause    *List
	ReturningList *List
	WithClause    *WithClause
}

func (n *UpdateStmt) Pos() int {
	return 0
}
