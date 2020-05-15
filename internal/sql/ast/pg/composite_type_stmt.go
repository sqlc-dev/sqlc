package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CompositeTypeStmt struct {
	Typevar    *RangeVar
	Coldeflist *ast.List
}

func (n *CompositeTypeStmt) Pos() int {
	return 0
}
