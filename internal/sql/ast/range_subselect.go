package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RangeSubselect struct {
	Lateral  bool
	Subquery ast.Node
	Alias    *Alias
}

func (n *RangeSubselect) Pos() int {
	return 0
}
