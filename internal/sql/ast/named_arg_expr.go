package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type NamedArgExpr struct {
	Xpr       Node
	Arg       Node
	Name      *string
	Argnumber int
	Location  int
}

func (n *NamedArgExpr) Pos() int {
	return n.Location
}
