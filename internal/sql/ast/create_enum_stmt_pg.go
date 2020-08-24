package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateEnumStmt struct {
	TypeName *List
	Vals     *List
}

func (n *CreateEnumStmt) Pos() int {
	return 0
}
