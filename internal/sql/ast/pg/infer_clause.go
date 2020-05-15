package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type InferClause struct {
	IndexElems  *ast.List
	WhereClause ast.Node
	Conname     *string
	Location    int
}

func (n *InferClause) Pos() int {
	return n.Location
}
