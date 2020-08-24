package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type IntoClause struct {
	Rel            *RangeVar
	ColNames       *List
	Options        *List
	OnCommit       OnCommitAction
	TableSpaceName *string
	ViewQuery      ast.Node
	SkipData       bool
}

func (n *IntoClause) Pos() int {
	return 0
}
