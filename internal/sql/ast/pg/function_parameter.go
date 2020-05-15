package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type FunctionParameter struct {
	Name    *string
	ArgType *TypeName
	Mode    FunctionParameterMode
	Defexpr ast.Node
}

func (n *FunctionParameter) Pos() int {
	return 0
}
