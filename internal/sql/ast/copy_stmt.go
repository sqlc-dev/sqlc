package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CopyStmt struct {
	Relation  *RangeVar
	Query     ast.Node
	Attlist   *List
	IsFrom    bool
	IsProgram bool
	Filename  *string
	Options   *List
}

func (n *CopyStmt) Pos() int {
	return 0
}
