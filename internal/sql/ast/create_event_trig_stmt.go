package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateEventTrigStmt struct {
	Trigname   *string
	Eventname  *string
	Whenclause *List
	Funcname   *List
}

func (n *CreateEventTrigStmt) Pos() int {
	return 0
}
