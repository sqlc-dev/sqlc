package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type NextValueExpr struct {
	Xpr    ast.Node
	Seqid  Oid
	TypeId Oid
}

func (n *NextValueExpr) Pos() int {
	return 0
}
