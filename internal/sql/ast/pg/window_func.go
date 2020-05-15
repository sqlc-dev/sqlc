package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type WindowFunc struct {
	Xpr         ast.Node
	Winfnoid    Oid
	Wintype     Oid
	Wincollid   Oid
	Inputcollid Oid
	Args        *ast.List
	Aggfilter   ast.Node
	Winref      Index
	Winstar     bool
	Winagg      bool
	Location    int
}

func (n *WindowFunc) Pos() int {
	return n.Location
}
