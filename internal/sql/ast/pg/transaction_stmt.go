package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type TransactionStmt struct {
	Kind    TransactionStmtKind
	Options *ast.List
	Gid     *string
}

func (n *TransactionStmt) Pos() int {
	return 0
}
