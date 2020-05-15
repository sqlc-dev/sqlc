package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateEventTrigStmt struct {
	Trigname   *string
	Eventname  *string
	Whenclause *ast.List
	Funcname   *ast.List
}

func (n *CreateEventTrigStmt) Pos() int {
	return 0
}
