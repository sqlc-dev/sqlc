package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type TransactionStmt struct {
	Kind    TransactionStmtKind
	Options *List
	Gid     *string
}

func (n *TransactionStmt) Pos() int {
	return 0
}
