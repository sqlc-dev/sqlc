package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type Alias struct {
	Aliasname *string
	Colnames  *ast.List
}

func (n *Alias) Pos() int {
	return 0
}
