package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RangeTableFunc struct {
	Lateral    bool
	Docexpr    ast.Node
	Rowexpr    ast.Node
	Namespaces *List
	Columns    *List
	Alias      *Alias
	Location   int
}

func (n *RangeTableFunc) Pos() int {
	return n.Location
}
