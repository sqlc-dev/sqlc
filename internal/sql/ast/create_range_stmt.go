package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateRangeStmt struct {
	TypeName *List
	Params   *List
}

func (n *CreateRangeStmt) Pos() int {
	return 0
}
