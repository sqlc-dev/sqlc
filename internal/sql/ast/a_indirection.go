package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type A_Indirection struct {
	Arg         ast.Node
	Indirection *List
}

func (n *A_Indirection) Pos() int {
	return 0
}
