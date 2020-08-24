package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RangeFunction struct {
	Lateral    bool
	Ordinality bool
	IsRowsfrom bool
	Functions  *List
	Alias      *Alias
	Coldeflist *List
}

func (n *RangeFunction) Pos() int {
	return 0
}
