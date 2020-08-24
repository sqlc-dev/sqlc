package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterFdwStmt struct {
	Fdwname     *string
	FuncOptions *List
	Options     *List
}

func (n *AlterFdwStmt) Pos() int {
	return 0
}
