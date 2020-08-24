package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreatePolicyStmt struct {
	PolicyName *string
	Table      *RangeVar
	CmdName    *string
	Permissive bool
	Roles      *ast.List
	Qual       ast.Node
	WithCheck  ast.Node
}

func (n *CreatePolicyStmt) Pos() int {
	return 0
}
