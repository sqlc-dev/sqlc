package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ObjectWithArgs struct {
	Objname         *ast.List
	Objargs         *ast.List
	ArgsUnspecified bool
}

func (n *ObjectWithArgs) Pos() int {
	return 0
}
