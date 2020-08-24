package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type WithClause struct {
	Ctes      *List
	Recursive bool
	Location  int
}

func (n *WithClause) Pos() int {
	return n.Location
}
