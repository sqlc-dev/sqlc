package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterPolicyStmt struct {
	PolicyName *string
	Table      *RangeVar
	Roles      *List
	Qual       Node
	WithCheck  Node
}

func (n *AlterPolicyStmt) Pos() int {
	return 0
}
