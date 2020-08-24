package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlternativeSubPlan struct {
	Xpr      ast.Node
	Subplans *List
}

func (n *AlternativeSubPlan) Pos() int {
	return 0
}
