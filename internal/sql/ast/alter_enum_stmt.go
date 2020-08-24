package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterEnumStmt struct {
	TypeName           *List
	OldVal             *string
	NewVal             *string
	NewValNeighbor     *string
	NewValIsAfter      bool
	SkipIfNewValExists bool
}

func (n *AlterEnumStmt) Pos() int {
	return 0
}
