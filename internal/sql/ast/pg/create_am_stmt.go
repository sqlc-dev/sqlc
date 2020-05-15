package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateAmStmt struct {
	Amname      *string
	HandlerName *ast.List
	Amtype      byte
}

func (n *CreateAmStmt) Pos() int {
	return 0
}
