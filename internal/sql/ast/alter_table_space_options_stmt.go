package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterTableSpaceOptionsStmt struct {
	Tablespacename *string
	Options        *ast.List
	IsReset        bool
}

func (n *AlterTableSpaceOptionsStmt) Pos() int {
	return 0
}
