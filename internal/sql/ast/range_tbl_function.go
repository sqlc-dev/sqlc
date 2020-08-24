package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RangeTblFunction struct {
	Funcexpr          ast.Node
	Funccolcount      int
	Funccolnames      *List
	Funccoltypes      *List
	Funccoltypmods    *List
	Funccolcollations *List
	Funcparams        []uint32
}

func (n *RangeTblFunction) Pos() int {
	return 0
}
