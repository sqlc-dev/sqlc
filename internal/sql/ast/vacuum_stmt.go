package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type VacuumStmt struct {
	Options  int
	Relation *RangeVar
	VaCols   *ast.List
}

func (n *VacuumStmt) Pos() int {
	return 0
}
