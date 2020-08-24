package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type DeleteStmt struct {
	Relation      *RangeVar
	UsingClause   *List
	WhereClause   ast.Node
	ReturningList *List
	WithClause    *WithClause
}

func (n *DeleteStmt) Pos() int {
	return 0
}
