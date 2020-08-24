package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type FieldStore struct {
	Xpr        ast.Node
	Arg        ast.Node
	Newvals    *ast.List
	Fieldnums  *ast.List
	Resulttype Oid
}

func (n *FieldStore) Pos() int {
	return 0
}
