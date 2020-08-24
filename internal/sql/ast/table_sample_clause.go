package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type TableSampleClause struct {
	Tsmhandler Oid
	Args       *List
	Repeatable ast.Node
}

func (n *TableSampleClause) Pos() int {
	return 0
}
