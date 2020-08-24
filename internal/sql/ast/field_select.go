package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type FieldSelect struct {
	Xpr          ast.Node
	Arg          ast.Node
	Fieldnum     AttrNumber
	Resulttype   Oid
	Resulttypmod int32
	Resultcollid Oid
}

func (n *FieldSelect) Pos() int {
	return 0
}
