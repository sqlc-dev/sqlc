package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AccessPriv struct {
	PrivName *string
	Cols     *ast.List
}

func (n *AccessPriv) Pos() int {
	return 0
}
