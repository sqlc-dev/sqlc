package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type InferenceElem struct {
	Xpr          ast.Node
	Expr         ast.Node
	Infercollid  Oid
	Inferopclass Oid
}

func (n *InferenceElem) Pos() int {
	return 0
}
