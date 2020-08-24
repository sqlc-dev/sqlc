package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type WindowClause struct {
	Name            *string
	Refname         *string
	PartitionClause *List
	OrderClause     *List
	FrameOptions    int
	StartOffset     ast.Node
	EndOffset       ast.Node
	Winref          Index
	CopiedOrder     bool
}

func (n *WindowClause) Pos() int {
	return 0
}
