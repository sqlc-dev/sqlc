package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type Aggref struct {
	Xpr           ast.Node
	Aggfnoid      Oid
	Aggtype       Oid
	Aggcollid     Oid
	Inputcollid   Oid
	Aggtranstype  Oid
	Aggargtypes   *ast.List
	Aggdirectargs *ast.List
	Args          *ast.List
	Aggorder      *ast.List
	Aggdistinct   *ast.List
	Aggfilter     ast.Node
	Aggstar       bool
	Aggvariadic   bool
	Aggkind       byte
	Agglevelsup   Index
	Aggsplit      AggSplit
	Location      int
}

func (n *Aggref) Pos() int {
	return n.Location
}
