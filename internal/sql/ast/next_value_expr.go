package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type NextValueExpr struct {
	Xpr    Node
	Seqid  Oid
	TypeId Oid
}

func (n *NextValueExpr) Pos() int {
	return 0
}
