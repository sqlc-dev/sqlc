package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type DefElem struct {
	Defnamespace *string
	Defname      *string
	Arg          ast.Node
	Defaction    DefElemAction
	Location     int
}

func (n *DefElem) Pos() int {
	return n.Location
}
