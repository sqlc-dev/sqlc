package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CoerceToDomainValue struct {
	Xpr       ast.Node
	TypeId    Oid
	TypeMod   int32
	Collation Oid
	Location  int
}

func (n *CoerceToDomainValue) Pos() int {
	return n.Location
}
