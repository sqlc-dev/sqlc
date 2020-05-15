package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterTableMoveAllStmt struct {
	OrigTablespacename *string
	Objtype            ObjectType
	Roles              *ast.List
	NewTablespacename  *string
	Nowait             bool
}

func (n *AlterTableMoveAllStmt) Pos() int {
	return 0
}
