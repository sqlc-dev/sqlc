package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateFunctionStmt struct {
	Replace    bool
	Funcname   *List
	Parameters *List
	ReturnType *TypeName
	Options    *List
	WithClause *List
}

func (n *CreateFunctionStmt) Pos() int {
	return 0
}
