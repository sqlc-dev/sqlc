package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type VacuumStmt struct {
	Options  int
	Relation *RangeVar
	VaCols   *List
}

func (n *VacuumStmt) Pos() int {
	return 0
}
