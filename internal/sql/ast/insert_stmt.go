package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type InsertStmt struct {
	Relation         *RangeVar
	Cols             *List
	SelectStmt       ast.Node
	OnConflictClause *OnConflictClause
	ReturningList    *List
	WithClause       *WithClause
	Override         OverridingKind
}

func (n *InsertStmt) Pos() int {
	return 0
}
