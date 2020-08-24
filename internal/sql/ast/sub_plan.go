package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type SubPlan struct {
	Xpr               ast.Node
	SubLinkType       SubLinkType
	Testexpr          ast.Node
	ParamIds          *List
	PlanId            int
	PlanName          *string
	FirstColType      Oid
	FirstColTypmod    int32
	FirstColCollation Oid
	UseHashTable      bool
	UnknownEqFalse    bool
	ParallelSafe      bool
	SetParam          *List
	ParParam          *List
	Args              *List
	StartupCost       Cost
	PerCallCost       Cost
}

func (n *SubPlan) Pos() int {
	return 0
}
