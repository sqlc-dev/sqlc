package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type TableSampleClause struct {
	Tsmhandler Oid
	Args       *List
	Repeatable Node
}

func (n *TableSampleClause) Pos() int {
	return 0
}
