package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterDatabaseStmt struct {
	Dbname  *string
	Options *List
}

func (n *AlterDatabaseStmt) Pos() int {
	return 0
}
