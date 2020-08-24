package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type InsertStmt struct {
	Relation         *RangeVar
	Cols             *ast.List
	SelectStmt       ast.Node
	OnConflictClause *OnConflictClause
	ReturningList    *ast.List
	WithClause       *WithClause
	Override         OverridingKind
}

func (n *InsertStmt) Pos() int {
	return 0
}
