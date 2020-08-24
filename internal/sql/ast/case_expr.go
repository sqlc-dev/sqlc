package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CaseExpr struct {
	Xpr        ast.Node
	Casetype   Oid
	Casecollid Oid
	Arg        ast.Node
	Args       *List
	Defresult  ast.Node
	Location   int
}

func (n *CaseExpr) Pos() int {
	return n.Location
}
