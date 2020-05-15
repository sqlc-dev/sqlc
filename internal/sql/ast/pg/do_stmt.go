package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type DoStmt struct {
	Args *ast.List
}

func (n *DoStmt) Pos() int {
	return 0
}
