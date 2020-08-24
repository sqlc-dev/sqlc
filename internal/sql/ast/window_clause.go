package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type WindowClause struct {
	Name            *string
	Refname         *string
	PartitionClause *ast.List
	OrderClause     *ast.List
	FrameOptions    int
	StartOffset     ast.Node
	EndOffset       ast.Node
	Winref          Index
	CopiedOrder     bool
}

func (n *WindowClause) Pos() int {
	return 0
}
