package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RangeSubselect struct {
	Lateral  bool
	Subquery Node
	Alias    *Alias
}

func (n *RangeSubselect) Pos() int {
	return 0
}
