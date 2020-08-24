package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type PrepareStmt struct {
	Name     *string
	Argtypes *List
	Query    Node
}

func (n *PrepareStmt) Pos() int {
	return 0
}
