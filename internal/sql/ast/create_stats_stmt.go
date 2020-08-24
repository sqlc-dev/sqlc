package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateStatsStmt struct {
	Defnames    *List
	StatTypes   *List
	Exprs       *List
	Relations   *List
	IfNotExists bool
}

func (n *CreateStatsStmt) Pos() int {
	return 0
}
