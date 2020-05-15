package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterTableStmt struct {
	Relation  *RangeVar
	Cmds      *ast.List
	Relkind   ObjectType
	MissingOk bool
}

func (n *AlterTableStmt) Pos() int {
	return 0
}
