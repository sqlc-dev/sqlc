package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterObjectDependsStmt struct {
	ObjectType ObjectType
	Relation   *RangeVar
	Object     ast.Node
	Extname    ast.Node
}

func (n *AlterObjectDependsStmt) Pos() int {
	return 0
}
