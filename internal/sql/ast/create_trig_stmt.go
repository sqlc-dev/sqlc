package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateTrigStmt struct {
	Trigname       *string
	Relation       *RangeVar
	Funcname       *List
	Args           *List
	Row            bool
	Timing         int16
	Events         int16
	Columns        *List
	WhenClause     ast.Node
	Isconstraint   bool
	TransitionRels *List
	Deferrable     bool
	Initdeferred   bool
	Constrrel      *RangeVar
}

func (n *CreateTrigStmt) Pos() int {
	return 0
}
