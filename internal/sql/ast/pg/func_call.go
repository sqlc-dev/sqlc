package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type FuncCall struct {
	Funcname       *ast.List
	Args           *ast.List
	AggOrder       *ast.List
	AggFilter      ast.Node
	AggWithinGroup bool
	AggStar        bool
	AggDistinct    bool
	FuncVariadic   bool
	Over           *WindowDef
	Location       int
}

func (n *FuncCall) Pos() int {
	return n.Location
}
