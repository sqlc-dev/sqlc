package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateFdwStmt struct {
	Fdwname     *string
	FuncOptions *List
	Options     *List
}

func (n *CreateFdwStmt) Pos() int {
	return 0
}
