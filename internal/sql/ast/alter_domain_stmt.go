package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterDomainStmt struct {
	Subtype   byte
	TypeName  *ast.List
	Name      *string
	Def       ast.Node
	Behavior  DropBehavior
	MissingOk bool
}

func (n *AlterDomainStmt) Pos() int {
	return 0
}
