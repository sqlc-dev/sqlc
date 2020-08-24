package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type Param struct {
	Xpr         ast.Node
	Paramkind   ParamKind
	Paramid     int
	Paramtype   Oid
	Paramtypmod int32
	Paramcollid Oid
	Location    int
}

func (n *Param) Pos() int {
	return n.Location
}
