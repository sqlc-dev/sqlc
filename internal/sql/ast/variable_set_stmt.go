package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type VariableSetStmt struct {
	Kind    VariableSetKind
	Name    *string
	Args    *List
	IsLocal bool
}

func (n *VariableSetStmt) Pos() int {
	return 0
}
