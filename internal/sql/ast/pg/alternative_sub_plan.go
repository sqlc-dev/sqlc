package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlternativeSubPlan struct {
	Xpr      ast.Node
	Subplans *ast.List
}

func (n *AlternativeSubPlan) Pos() int {
	return 0
}
