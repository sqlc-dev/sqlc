package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlternativeSubPlan struct {
	Xpr      Node
	Subplans *List
}

func (n *AlternativeSubPlan) Pos() int {
	return 0
}
