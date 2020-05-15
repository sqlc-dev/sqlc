package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterDatabaseStmt struct {
	Dbname  *string
	Options *ast.List
}

func (n *AlterDatabaseStmt) Pos() int {
	return 0
}
