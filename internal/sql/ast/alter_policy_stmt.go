package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterPolicyStmt struct {
	PolicyName *string
	Table      *RangeVar
	Roles      *List
	Qual       ast.Node
	WithCheck  ast.Node
}

func (n *AlterPolicyStmt) Pos() int {
	return 0
}
