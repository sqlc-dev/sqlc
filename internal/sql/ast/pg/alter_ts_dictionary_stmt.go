package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterTSDictionaryStmt struct {
	Dictname *ast.List
	Options  *ast.List
}

func (n *AlterTSDictionaryStmt) Pos() int {
	return 0
}
