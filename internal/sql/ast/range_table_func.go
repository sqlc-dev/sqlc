package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RangeTableFunc struct {
	Lateral    bool
	Docexpr    ast.Node
	Rowexpr    ast.Node
	Namespaces *ast.List
	Columns    *ast.List
	Alias      *Alias
	Location   int
}

func (n *RangeTableFunc) Pos() int {
	return n.Location
}
