package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateTrigStmt struct {
	Trigname       *string
	Relation       *RangeVar
	Funcname       *ast.List
	Args           *ast.List
	Row            bool
	Timing         int16
	Events         int16
	Columns        *ast.List
	WhenClause     ast.Node
	Isconstraint   bool
	TransitionRels *ast.List
	Deferrable     bool
	Initdeferred   bool
	Constrrel      *RangeVar
}

func (n *CreateTrigStmt) Pos() int {
	return 0
}
