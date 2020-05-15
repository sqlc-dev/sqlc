package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterTableCmd struct {
	Subtype   AlterTableType
	Name      *string
	Newowner  *RoleSpec
	Def       ast.Node
	Behavior  DropBehavior
	MissingOk bool
}

func (n *AlterTableCmd) Pos() int {
	return 0
}
