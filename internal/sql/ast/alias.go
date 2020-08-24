package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type Alias struct {
	Aliasname *string
	Colnames  *List
}

func (n *Alias) Pos() int {
	return 0
}
