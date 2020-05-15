package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type WindowDef struct {
	Name            *string
	Refname         *string
	PartitionClause *ast.List
	OrderClause     *ast.List
	FrameOptions    int
	StartOffset     ast.Node
	EndOffset       ast.Node
	Location        int
}

func (n *WindowDef) Pos() int {
	return n.Location
}
