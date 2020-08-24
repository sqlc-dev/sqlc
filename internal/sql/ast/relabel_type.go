package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RelabelType struct {
	Xpr           ast.Node
	Arg           ast.Node
	Resulttype    Oid
	Resulttypmod  int32
	Resultcollid  Oid
	Relabelformat CoercionForm
	Location      int
}

func (n *RelabelType) Pos() int {
	return n.Location
}
