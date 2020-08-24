package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterCollationStmt struct {
	Collname *List
}

func (n *AlterCollationStmt) Pos() int {
	return 0
}
