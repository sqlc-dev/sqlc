package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type IndexElem struct {
	Name          *string
	Expr          ast.Node
	Indexcolname  *string
	Collation     *ast.List
	Opclass       *ast.List
	Ordering      SortByDir
	NullsOrdering SortByNulls
}

func (n *IndexElem) Pos() int {
	return 0
}
