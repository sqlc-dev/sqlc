package ast

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
	Aggargtypes   *List
	Aggdirectargs *List
	Args          *List
	Aggorder      *List
	Aggdistinct   *List
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
