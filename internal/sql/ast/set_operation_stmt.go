package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type SetOperationStmt struct {
	Op            SetOperation
	All           bool
	Larg          ast.Node
	Rarg          ast.Node
	ColTypes      *ast.List
	ColTypmods    *ast.List
	ColCollations *ast.List
	GroupClauses  *ast.List
}

func (n *SetOperationStmt) Pos() int {
	return 0
}
