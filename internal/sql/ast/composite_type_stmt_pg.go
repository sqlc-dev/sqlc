package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CompositeTypeStmt struct {
	Typevar    *RangeVar
	Coldeflist *List
}

func (n *CompositeTypeStmt) Pos() int {
	return 0
}
