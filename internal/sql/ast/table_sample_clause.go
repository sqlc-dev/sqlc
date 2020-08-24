package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type TableSampleClause struct {
	Tsmhandler Oid
	Args       *ast.List
	Repeatable ast.Node
}

func (n *TableSampleClause) Pos() int {
	return 0
}
