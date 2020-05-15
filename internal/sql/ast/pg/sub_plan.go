package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type SubPlan struct {
	Xpr               ast.Node
	SubLinkType       SubLinkType
	Testexpr          ast.Node
	ParamIds          *ast.List
	PlanId            int
	PlanName          *string
	FirstColType      Oid
	FirstColTypmod    int32
	FirstColCollation Oid
	UseHashTable      bool
	UnknownEqFalse    bool
	ParallelSafe      bool
	SetParam          *ast.List
	ParParam          *ast.List
	Args              *ast.List
	StartupCost       Cost
	PerCallCost       Cost
}

func (n *SubPlan) Pos() int {
	return 0
}
