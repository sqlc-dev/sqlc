package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type SetOperationStmt struct {
	Op            SetOperation
	All           bool
	Larg          ast.Node
	Rarg          ast.Node
	ColTypes      *List
	ColTypmods    *List
	ColCollations *List
	GroupClauses  *List
}

func (n *SetOperationStmt) Pos() int {
	return 0
}
