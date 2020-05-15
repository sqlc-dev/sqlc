package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CurrentOfExpr struct {
	Xpr         ast.Node
	Cvarno      Index
	CursorName  *string
	CursorParam int
}

func (n *CurrentOfExpr) Pos() int {
	return 0
}
